package search

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	SearchResponseError   = 0
	SearchResponsePage    = 1
	SearchResponseCaptcha = 2
)

type SearchResponse struct {
	Type       int
	Status     int
	Captcha    *CaptchaPage
	SearchPage *SearchPage
}

func Search(searchTerm string, offset int) (response SearchResponse, err error) {
	searchUrl := getSearchUrl(searchTerm, offset)
	document, err, status := getDocument(searchUrl)

	if err != nil {
		return SearchResponse{Type: SearchResponseError, Status: 0}, err
	}

	if status != http.StatusOK {
		if status == http.StatusTooManyRequests {
			captchaPage := createCaptchaPage(searchTerm)
			return SearchResponse{Type: SearchResponseCaptcha, Captcha: &captchaPage, Status: status}, err
		}
		return SearchResponse{Type: SearchResponseError, Status: status}, nil
	}

	if len(searchTerm) == 0 {
		searchPage := createEmptySearchPage()
		return SearchResponse{Type: SearchResponsePage, SearchPage: &searchPage, Status: status}, nil
	} else {
		searchPage, err := parseSearchPage(document, searchTerm, offset)

		if err != nil {
			return SearchResponse{Type: SearchResponseError, Status: status}, err
		}

		return SearchResponse{Type: SearchResponsePage, SearchPage: searchPage, Status: status}, nil
	}
}

func getSearchUrl(searchTerm string, start int) string {
	searchUrl, _ := url.Parse("https://google.com/search")
	escapedTerm := url.QueryEscape(searchTerm)
	query := searchUrl.Query()

	query.Add("q", escapedTerm)

	if start > 0 {
		query.Add("start", strconv.Itoa(start))
	}

	searchUrl.RawQuery = query.Encode()
	return searchUrl.String()
}

func getDocument(url string) (document *goquery.Document, err error, status int) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err, 0
	}

	req.Header = http.Header{
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8"},
		"Accept-Language": {"en-US,en;q=0.8"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err, res.StatusCode
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err, res.StatusCode
	}

	return doc, nil, res.StatusCode
}
