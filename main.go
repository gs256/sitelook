package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

func apiSearchRoute(c *gin.Context) {
	searchTerm := c.Query("q")
	startQuery := c.Query("start")

	if len(searchTerm) == 0 {
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.WriteString("empty search term")
		return
	}

	start, _ := strconv.Atoi(startQuery)
	searchPage, err := parseSearchPage(searchTerm, start)

	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, searchPage)
}

type SearchResultContext struct {
	Url         string
	Title       string
	UrlTitle    string
	Description string
}

type PageLinkContext struct {
	PageNumber int
	PageUrl    string
	IsCurrent  bool
}

type PaginationContext struct {
	Visible            bool
	PageLinks          []PageLinkContext
	PreviousUrl        string
	PreviousLinkActive bool
	NextUrl            string
	NextLinkActive     bool
}

type SearchNavigationContext struct {
	AllSearchHref   string
	ImageSearchHref string
	VideoSearchHref string
}

type SearchCorrectionContext struct {
	Present           bool
	Title             string
	CorrectSearchTerm string
	CorrectionHref    string
}

type SearchPageContext struct {
	SearchTerm       string
	SearchResults    []SearchResultContext
	Pagination       PaginationContext
	Navigation       SearchNavigationContext
	SearchCorrection SearchCorrectionContext
}

func createSearchCorrectionContext(searchCorrection SearchCorrection, currentUrl *url.URL) SearchCorrectionContext {
	query := currentUrl.Query()
	query.Set("q", searchCorrection.CorrectSearchTerm)
	href := createHref(currentUrl, query)

	return SearchCorrectionContext{
		Present:           searchCorrection.Present,
		Title:             searchCorrection.Title,
		CorrectSearchTerm: searchCorrection.CorrectSearchTerm,
		CorrectionHref:    href,
	}
}

func createSearchResultContext(searchResult SearchResult) SearchResultContext {
	u, _ := url.Parse(searchResult.Url)
	host := u.Host
	path := u.Path
	host = strings.TrimPrefix(host, "www.")
	urlTitle := strings.TrimRight(host+path, "/")

	return SearchResultContext{
		Url:         searchResult.Url,
		Title:       searchResult.Title,
		UrlTitle:    urlTitle,
		Description: searchResult.Description,
	}
}

func createPageLinkContext(pageLink PageLink, pageUrl string) PageLinkContext {
	return PageLinkContext{
		PageNumber: pageLink.PageNumber,
		PageUrl:    pageUrl,
		IsCurrent:  pageLink.IsCurrent,
	}
}

func createHref(url *url.URL, query url.Values) string {
	return url.Path + "?" + query.Encode()
}

func createPaginationContext(pagination Pagination, currentUrl *url.URL) PaginationContext {
	pageLinks := make([]PageLinkContext, len(pagination.PageLinks))
	query := currentUrl.Query()

	for i := 0; i < len(pagination.PageLinks); i++ {
		link := pagination.PageLinks[i]
		query.Set("start", strconv.Itoa(link.Offset))
		pageLinks[i] = createPageLinkContext(link, createHref(currentUrl, query))
	}

	query.Set("start", strconv.Itoa(pagination.PreviousOffset))
	previousUrl := createHref(currentUrl, query)
	spew.Dump(query, previousUrl)
	query.Set("start", strconv.Itoa(pagination.NextOffset))
	nextUrl := createHref(currentUrl, query)
	spew.Dump(query, nextUrl)

	return PaginationContext{
		Visible:            len(pageLinks) > 0,
		PageLinks:          pageLinks,
		PreviousUrl:        previousUrl,
		PreviousLinkActive: len(pageLinks) > 0 && !pageLinks[0].IsCurrent,
		NextUrl:            nextUrl,
		NextLinkActive:     len(pageLinks) > 0 && !pageLinks[len(pageLinks)-1].IsCurrent,
	}
}

func createNavigationContext(currentUrl *url.URL) SearchNavigationContext {
	query := currentUrl.Query()
	query.Del("tbm")
	allSearchHref := createHref(currentUrl, query)
	query.Set("tbm", "isch")
	imageSearchHref := createHref(currentUrl, query)
	query.Set("tbm", "vid")
	videoSearchHref := createHref(currentUrl, query)

	return SearchNavigationContext{
		AllSearchHref:   allSearchHref,
		ImageSearchHref: imageSearchHref,
		VideoSearchHref: videoSearchHref,
	}
}

func createSearchPageContext(searchPage SearchPage, currentUrl *url.URL) SearchPageContext {
	searchResults := make([]SearchResultContext, len(searchPage.SearchResults))

	for i := 0; i < len(searchPage.SearchResults); i++ {
		searchResults[i] = createSearchResultContext(searchPage.SearchResults[i])
	}

	return SearchPageContext{
		SearchTerm:       searchPage.SearchTerm,
		SearchResults:    searchResults,
		Pagination:       createPaginationContext(searchPage.Pagination, currentUrl),
		Navigation:       createNavigationContext(currentUrl),
		SearchCorrection: createSearchCorrectionContext(searchPage.SearchCorrection, currentUrl),
	}
}

func createEmptySearchPageContext() SearchPageContext {
	return SearchPageContext{}
}

func homeRoute(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/search")
}

func searchRoute(c *gin.Context) {
	searchTerm := c.Query("q")
	searchType := c.Query("tbm")
	if len(searchType) > 0 {
		c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("https://google.com/search?q=%s&tbm=%s", searchTerm, searchType))
		return
	}

	startQuery := c.Query("start")
	start, _ := strconv.Atoi(startQuery)

	var searchPageContext SearchPageContext

	if len(searchTerm) == 0 {
		searchPageContext = createEmptySearchPageContext()
	} else {
		searchPage, err := parseSearchPage(searchTerm, start)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}

		currentUrl := c.Request.URL
		searchPageContext = createSearchPageContext(*searchPage, currentUrl)

	}

	file, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(file)
	defer file.Close()
	spew.Fdump(log.Writer(), searchPageContext)

	c.HTML(http.StatusOK, "search-page.html", searchPageContext)
}

func main() {
	engine := gin.Default()
	engine.GET("/", homeRoute)
	engine.GET("/api/search", apiSearchRoute)
	engine.GET("/search", searchRoute)
	engine.Static("./static", "./static/")
	engine.LoadHTMLGlob("templates/*")
	engine.Run()
}
