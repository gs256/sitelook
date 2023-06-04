package search

import (
	"errors"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

type Pagination struct {
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
	Pagination       Pagination
	SearchCorrection SearchCorrection
}

type CaptchaPage struct {
	SearchTerm string
}

func parsePagination(document *goquery.Document) (Pagination, error) {
	paginationDiv := findSingle(document.Selection, "div[role=\"navigation\"] table[role=\"presentation\"]")

	if selectionEmpty(paginationDiv) {
		return Pagination{}, errors.New("pagination container not found")
	}

	pageLinks := []PageLink{}
	previousPageOffset := 0
	nextPageOffset := 0
	paginationDiv.Find("td").Each(func(i int, td *goquery.Selection) {
		if _, exists := td.Attr("role"); exists {
			if i == 0 {
				previousPageOffset, _ = getOffsetFromSelection(td)
			} else {
				nextPageOffset, _ = getOffsetFromSelection(td)
			}
		} else {
			offset, hrefExists := getOffsetFromSelection(td)
			hasSpanInside := false
			if !hrefExists {
				hasSpanInside = hasInside(td, "span")
			}
			number, err := strconv.Atoi(strings.TrimSpace(td.Text()))
			if err == nil && (hrefExists || hasSpanInside) {
				pageLink := PageLink{
					PageNumber: number,
					Offset:     offset,
					IsCurrent:  hasSpanInside,
				}
				pageLinks = append(pageLinks, pageLink)
			} else {
				html, _ := td.Html()
				log.Printf("error parsing pagination link: `%s`", html)
			}
		}
	})

	return Pagination{
		PageLinks:      pageLinks,
		PreviousOffset: previousPageOffset,
		NextOffset:     nextPageOffset,
	}, nil
}

func parseSearchResults(document *goquery.Document) []SearchResult {
	results := []SearchResult{}
	searchDiv := findSingle(document.Selection, "#search")

	searchDiv.Find(".g > div").Each(func(i int, searchItem *goquery.Selection) {
		// if the item has nested results it will force them to be parsed individually
		if hasInside(searchItem, ".g") {
			return
		}

		titleElement := findSingle(searchItem, "h3")
		if selectionEmpty(titleElement) {
			return
		}

		title := titleElement.Text()
		url, _ := findSingle(searchItem, "a").Attr("href")

		description := ""
		descriptionElement := findSingle(searchItem, "div[data-sncf=\"1\"]")

		if !selectionEmpty(descriptionElement) {
			description = descriptionElement.Text()
		} else {
			spans := searchItem.Find("span").Last()
			if spans.Length() > 0 {
				description = spans.Text()
			}
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
	correctionContainer := findSingle(document.Selection, "#taw")

	correction := SearchCorrection{
		Present: false,
	}

	if correctionContainer.Length() != 0 {
		correction.Present = true
		correction.Title = findSingle(correctionContainer, "p > span").Text()
		correctionHref, _ := findSingle(correctionContainer, "p > a").Attr("href")
		correctionUrl, _ := url.Parse(correctionHref)
		correctionSearch := correctionUrl.Query().Get("q")
		correction.CorrectSearchTerm = correctionSearch
	}

	return correction
}

func parseSearchInput(document *goquery.Document) string {
	searchInput := findSingle(document.Selection, "textarea")
	return searchInput.Text()
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
