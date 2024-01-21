package search

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	PaginationTypeMultiPage  = "MultiPagePagination"  // default google search pagination
	PaginationTypeSinglePage = "SinglePagePagination" // non-js image search pagination
)

type SearchResult struct {
	Url         string
	Title       string
	Description string
}

type PageLink struct {
	PageNumber int
	Offset     int
	IsCurrent  bool
}

type MultiPagePagination struct {
	Type           string
	PageLinks      []PageLink
	PreviousOffset int
	NextOffset     int
}

type SearchCorrection struct {
	Present           bool
	Title             string
	CorrectSearchTerm string
}

type SearchPage struct {
	SearchTerm       string
	SearchResults    []SearchResult
	Pagination       SinglePagePagination
	SearchCorrection SearchCorrection
}

type CaptchaPage struct {
	SearchTerm string
}

func parsePagination(document *goquery.Document) (SinglePagePagination, error) {
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

func parseSearchResults(document *goquery.Document) []SearchResult {
	results := []SearchResult{}

	document.Find(".fP1Qef").Each(func(i int, searchItem *goquery.Selection) {
		titleElement := findSingle(searchItem, "h3")
		if selectionEmpty(titleElement) {
			return
		}

		title := titleElement.Text()
		url, _ := findSingle(searchItem, "a").Attr("href")
		url = hrefFromQuery(url)

		description := ""
		descriptionElement := findSingle(searchItem, ".BNeawe.s3v9rd.AP7Wnd")

		if !selectionEmpty(descriptionElement) {
			description = descriptionElement.Text()
		}

		results = append(results, SearchResult{
			Url:         url,
			Title:       title,
			Description: description,
		})
	})

	return results
}

func parseSearchCorrection(document *goquery.Document) SearchCorrection {
	correctionContainer := findSingle(document.Selection, "#scc")

	correction := SearchCorrection{
		Present: false,
	}

	if correctionContainer.Length() != 0 {
		correction.Present = true
		title := findSingle(correctionContainer, ".EE3Upf").Text()
		correction.Title = strings.Split(title, ":")[0]
		correctionHref, _ := findSingle(correctionContainer, "a").Attr("href")
		correctionUrl, _ := url.Parse(correctionHref)
		correctionSearch := correctionUrl.Query().Get("q")
		correction.CorrectSearchTerm = correctionSearch
	}

	return correction
}

func parseSearchInput(document *goquery.Document) string {
	searchInput := findSingle(document.Selection, "input[name=\"q\"]")
	searchText := searchInput.AttrOr("value", "")
	return searchText
}

func parseSearchPage(document *goquery.Document, start int) (*SearchPage, error) {
	searchInput := parseSearchInput(document)
	searchResults := parseSearchResults(document)
	searchCorrection := parseSearchCorrection(document)
	pagination, err := parsePagination(document)

	if err != nil {
		log.Printf("pagination error: %s", err)
	}

	searchPage := SearchPage{
		SearchTerm:       searchInput,
		SearchResults:    searchResults,
		Pagination:       pagination,
		SearchCorrection: searchCorrection,
	}

	return &searchPage, nil
}

func getOffsetFromSelection(selection *goquery.Selection) (offset int, isSet bool) {
	href, exists := findSingle(selection, "a").Attr("href")

	if !exists {
		return 0, false
	}

	return getOffsetFromHref(href)
}

func getOffsetFromLink(linkElement *goquery.Selection) (offset int, isSet bool) {
	href, set := linkElement.Attr("href")

	if !set {
		return 0, false
	}

	return getOffsetFromHref(href)
}

func getOffsetFromHref(href string) (offset int, isSet bool) {
	hrefUrl, err := url.Parse(href)
	if err != nil {
		return 0, false
	}

	params, err := url.ParseQuery(hrefUrl.RawQuery)
	if err != nil {
		return 0, false
	}

	offsetParam := params.Get("start")
	if len(offsetParam) == 0 {
		return 0, false
	}

	offsetInt, err := strconv.Atoi(offsetParam)
	if err != nil {
		return 0, false

	}

	return offsetInt, true
}
