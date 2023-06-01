package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	Url         string
	Title       string
	Description string
}

func parse(url string) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8"},
		"Accept-Language": {"en-US,en;q=0.8"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d\n", len(doc.Text()))

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
			// fmt.Printf("%d. %s (%s) - %s\n", i, url, title, description)
		}
	})

	searchInput := doc.Find("textarea").First()
	fmt.Printf("Search term: %s\n", searchInput.Text())
}

func main() {
	parse(getSearchUrl("how long is an hour"))
}

func getSearchUrl(searchTerm string) string {
	escapedTerm := url.QueryEscape(searchTerm)
	return fmt.Sprintf("https://google.com/search?q=%s", escapedTerm)
}
