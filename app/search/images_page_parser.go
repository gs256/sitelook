package search

import (
	"errors"
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

type ImagesPage struct {
	SearchTerm   string
	ImageResults []ImageResult
}

func hrefFromQuery(url_ string) string {
	u, err := url.Parse(url_)
	if err != nil {
		return ""
	}
	query := u.Query().Get("q")
	return query
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

	return ImagesPage{
		SearchTerm:   searchTerm,
		ImageResults: imageResults,
	}, nil
}
