package search

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type ImageResult struct {
	Title         string
	UrlTitle      string
	ImageSrc      string
	TitleLinkHref string
	ImageLinkHref string
}

type SinglePagePagination struct {
	PreviousLinkPresent bool
	PreviousOffset      int
	NextLinkPresent     bool
	NextOffset          int
	CurrentTitle        string
}

type ImagesPage struct {
	SearchTerm   string
	ImageResults []ImageResult
	Pagination   SinglePagePagination
}

func hrefFromQuery(url_ string) string {
	u, err := url.Parse(url_)
	if err != nil {
		return ""
	}
	query := u.Query().Get("q")
	return query
}

func parseImagePagePagination(document *goquery.Document) (SinglePagePagination, error) {
	paginationTable := findSingle(document.Selection, "body > table")

	if selectionEmpty(paginationTable) {
		return SinglePagePagination{}, errors.New("pagination not found")
	}

	tds := paginationTable.Find("td")

	pagination := SinglePagePagination{
		PreviousLinkPresent: false,
		PreviousOffset:      0,
		NextLinkPresent:     false,
		NextOffset:          0,
		CurrentTitle:        "",
	}

	if tds.Length() == 0 {
		return SinglePagePagination{}, errors.New("pagination links not found")
	} else {
		containers := make([]*goquery.Selection, tds.Length())
		tds.Each(func(i int, td *goquery.Selection) {
			containers[i] = td
		})

		if len(containers) == 1 {
			pagination.NextLinkPresent = true
			pagination.NextOffset, _ = getOffsetFromHref(findSingle(tds, "a").AttrOr("href", "#"))
		} else if len(containers) == 5 {
			previousLink := findSingle(containers[1], "a")
			pagination.PreviousLinkPresent = !selectionEmpty(previousLink)
			pagination.PreviousOffset, _ = getOffsetFromHref(previousLink.AttrOr("href", "#"))
			pagination.CurrentTitle = containers[2].Text()
			nextLink := findSingle(containers[3], "a")
			pagination.NextLinkPresent = !selectionEmpty(nextLink)
			nextOffset, isSet := getOffsetFromHref(nextLink.AttrOr("href", "#"))
			if isSet {
				pagination.NextOffset = nextOffset
			} else {
				pagination.NextLinkPresent = false
			}
		} else {
			return SinglePagePagination{}, errors.New(fmt.Sprintf("pagination has %d <td> elements (1 or 5 expected)", len(containers)))
		}
	}

	return pagination, nil
}

func selectionToArray(selection *goquery.Selection) []*goquery.Selection {
	elements := make([]*goquery.Selection, selection.Length())
	selection.Each(func(i int, element *goquery.Selection) {
		elements[i] = element
	})
	return elements
}

func parseVideoPagePagination(document *goquery.Document) (SinglePagePagination, error) {
	footer := findSingle(document.Selection, "footer > div")
	paginationDiv := findSingle(footer, "div > a").Parent()

	if selectionEmpty(paginationDiv) {
		return SinglePagePagination{}, errors.New("pagination not found")
	}

	pagination := SinglePagePagination{
		PreviousLinkPresent: false,
		PreviousOffset:      0,
		NextLinkPresent:     false,
		NextOffset:          0,
		CurrentTitle:        "",
	}

	links := selectionToArray(paginationDiv.Find("a"))
	span := findSingle(paginationDiv, "div > span")

	if len(links) == 0 {
		return pagination, errors.New("pagination links not found")
	} else if len(links) == 1 {
		if !selectionEmpty(span) {
			// Second page (last)
			pagination.PreviousLinkPresent = true
			pagination.PreviousOffset, _ = getOffsetFromLink(links[0])
			pagination.CurrentTitle = span.Text()
		} else {
			// First page
			pagination.NextLinkPresent = true
			pagination.NextOffset, _ = getOffsetFromLink(links[0])
		}
	} else if len(links) == 2 {
		paginationChildren := paginationDiv.Children()

		if paginationChildren.Last().Is("span") {
			// Last page (3+)
			pagination.PreviousLinkPresent = true
			pagination.PreviousOffset, _ = getOffsetFromLink(links[1])
			pagination.CurrentTitle = span.Text()
		} else {
			// Second page (if not last)
			pagination.PreviousLinkPresent = true
			pagination.PreviousOffset, _ = getOffsetFromLink(links[0])
			pagination.NextLinkPresent = true
			pagination.NextOffset, _ = getOffsetFromLink(links[1])
			pagination.CurrentTitle = span.Text()
		}
	} else if len(links) == 3 {
		// Any page between second and last (but not second and not last)
		pagination.PreviousLinkPresent = true
		pagination.PreviousOffset, _ = getOffsetFromLink(links[1])
		pagination.NextLinkPresent = true
		pagination.NextOffset, _ = getOffsetFromLink(links[2])
		pagination.CurrentTitle = span.Text()
	} else {
		return pagination, errors.New(fmt.Sprintf("pagination has %d links (1-3 expected)", len(links)))
	}

	return pagination, nil
}

func parseImagesPage(document *goquery.Document) (ImagesPage, error) {
	searchInput := findSingle(document.Selection, "input[name=\"q\"]")

	if selectionEmpty(searchInput) {
		return ImagesPage{}, errors.New("search input not found")
	}

	searchTerm := searchInput.AttrOr("value", "")

	imageResults := make([]ImageResult, 0)

	document.Find("tbody").Each(func(i int, tbody *goquery.Selection) {
		image := tbody.Find("img")
		if image.Length() != 1 {
			return
		}

		imageSrc := image.AttrOr("src", "#")

		links := tbody.Find("a")
		if links.Length() != 2 {
			log.Fatalf("image result element has %d links instead of 2", links.Length())
			return
		}

		// extracting image url from `<a href="/url?q=[image url]">...</a>`
		imageLinkHref := hrefFromQuery(links.First().AttrOr("href", "#"))
		titleLinkHref := hrefFromQuery(links.Last().AttrOr("href", "#"))

		spans := tbody.Find("a span > span")
		if spans.Length() != 2 {
			log.Fatalf("image element has %d spans instead of 2", spans.Length())
			return
		}

		title := spans.First().Text()
		urlTitle := spans.Last().Text()

		imageResults = append(imageResults, ImageResult{
			Title:         title,
			UrlTitle:      urlTitle,
			ImageSrc:      imageSrc,
			TitleLinkHref: titleLinkHref,
			ImageLinkHref: imageLinkHref,
		})
	})

	if len(imageResults) == 0 {
		return ImagesPage{
			SearchTerm:   searchTerm,
			ImageResults: imageResults,
		}, errors.New("page has no images or an error occured while parsing images")
	}

	pagination, err := parseImagePagePagination(document)
	if err != nil {
		log.Println(err)
	}

	return ImagesPage{
		SearchTerm:   searchTerm,
		ImageResults: imageResults,
		Pagination:   pagination,
	}, nil
}
