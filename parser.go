package main

import (
	"errors"
	"fmt"
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

func selectionEmpty(selection *goquery.Selection) bool {
	return selection.Length() == 0
}

func parsePagination(document *goquery.Document) (Pagination, error) {
	paginationDiv := document.Find("div[role=\"navigation\"] table[role=\"presentation\"]").First()

	if selectionEmpty(paginationDiv) {
		return Pagination{}, errors.New("pagination container not found")
	}

	pageLinks := []PageLink{}
	previousPageOffset := 0
	nextPageOffset := 0
	paginationDiv.Find("td").Each(func(i int, s *goquery.Selection) {
		if _, exists := s.Attr("role"); exists {
			if i == 0 {
				previousPageOffset, _ = getOffsetFromSelection(s)
			} else {
				nextPageOffset, _ = getOffsetFromSelection(s)
			}
		} else {
			offset, hrefExists := getOffsetFromSelection(s)
			hasSpanInside := false
			if !hrefExists {
				hasSpanInside = s.Find("span").Length() > 0
			}
			number, err := strconv.Atoi(strings.TrimSpace(s.Text()))
			if err == nil && (hrefExists || hasSpanInside) {
				pageLink := PageLink{
					PageNumber: number,
					Offset:     offset,
					IsCurrent:  hasSpanInside,
				}
				pageLinks = append(pageLinks, pageLink)
			} else {
				html, _ := s.Html()
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
	searchDiv := document.Find("#search").First()

	searchDiv.Find(".g > div").Each(func(i int, searchItem *goquery.Selection) {
		// if the item has nested results it will force them to be parsed individually
		if !selectionEmpty(searchItem.Find(".g")) {
			return
		}

		titleElement := searchItem.Find("h3").First()
		if selectionEmpty(titleElement) {
			return
		}

		title := titleElement.Text()
		url, _ := searchItem.Find("a").First().Attr("href")

		description := ""
		descriptionElement := searchItem.Find("div[data-sncf=\"1\"]").First()

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
	correctionContainer := document.Find("#taw").First()

	correction := SearchCorrection{
		Present: false,
	}

	if correctionContainer.Length() != 0 {
		correction.Present = true
		correction.Title = correctionContainer.Find("p > span").First().Text()
		correctionHref, _ := correctionContainer.Find("p > a").First().Attr("href")
		correctionUrl, _ := url.Parse(correctionHref)
		correctionSearch := correctionUrl.Query().Get("q")
		correction.CorrectSearchTerm = correctionSearch
	}

	return correction
}

func parseSearchInput(document *goquery.Document) string {
	searchInput := document.Find("textarea").First()
	return searchInput.Text()
}

func parseSearchPage(document *goquery.Document, searchTerm string, start int) (*SearchPage, error) {
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

func getOffsetFromSelection(s *goquery.Selection) (offset int, isSet bool) {
	href, exists := s.Find("a").First().Attr("href")
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

func getSearchUrl(searchTerm string, start int) string {
	escapedTerm := url.QueryEscape(searchTerm)
	return fmt.Sprintf("https://google.com/search?q=%s&start=%d", escapedTerm, start)
}
