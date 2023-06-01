package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
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

type Paginaiton struct {
	PageLinks      []PageLink
	NextOffset     int
	PreviousOffset int
}

type SearchPageContext struct {
	SearchTerm    string
	SearchResults []SearchResult
	Pagination    Paginaiton
}

func getDocument(url string) (*goquery.Document, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8"},
		"Accept-Language": {"en-US,en;q=0.8"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func parseSearchPage(searchTerm string, start int) (*SearchPageContext, error) {
	url := getSearchUrl(searchTerm, start)
	doc, err := getDocument(url)

	if err != nil {
		return nil, err
	}

	paginationDiv := doc.Find("div[role=\"navigation\"] table[role=\"presentation\"]").First()
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

	results := []SearchResult{}

	searchDiv := doc.Find("#search").First()
	searchDiv.Find("div[eid] > div").Each(func(i int, s *goquery.Selection) {
		url, hasUrl := s.Find("div[data-snf] a").First().Attr("href")
		title := s.Find("div[data-snhf=\"0\"] h3").First().Text()
		description := s.Find("div[data-sncf=\"1\"] div:last-of-type").First().Text()

		if hasUrl && len(title) > 0 {
			results = append(results, SearchResult{
				Url:         url,
				Title:       title,
				Description: description,
			})
		}
	})

	searchInput := doc.Find("textarea").First()

	context := SearchPageContext{
		SearchTerm:    searchInput.Text(),
		SearchResults: results,
		Pagination: Paginaiton{
			PageLinks:      pageLinks,
			NextOffset:     nextPageOffset,
			PreviousOffset: previousPageOffset,
		},
	}

	return &context, nil
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

func main() {
	searchPageContext, err := parseSearchPage("tesr", 80)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(searchPageContext)
}
