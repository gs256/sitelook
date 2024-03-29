package search

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

// func apiSearchRoute(c *gin.Context) {
// 	searchTerm := c.Query("q")
// 	startQuery := c.Query("start")

// 	if len(searchTerm) == 0 {
// 		c.Writer.WriteHeader(http.StatusBadRequest)
// 		c.Writer.WriteString("empty search term")
// 		return
// 	}

// 	start, _ := strconv.Atoi(startQuery)
// 	searchPage, err := parseSearchPage(searchTerm, start)

// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusInternalServerError)
// 		log.Fatal(err)
// 	}

// 	c.JSON(http.StatusOK, searchPage)
// }

func createHref(url *url.URL, query url.Values) string {
	return url.Path + "?" + query.Encode()
}

type SearchQueryParams struct {
	Type              string
	Start             int
	SearchLanguage    string
	InterfaceLanguage string
}

func createSearchQueryParams(context *gin.Context) SearchQueryParams {
	searchType := context.Query("tbm")
	startQuery := context.Query("start")
	lrQuery := context.Query("lr")
	hlQuery := context.Query("hl")
	start, _ := strconv.Atoi(startQuery)

	return SearchQueryParams{
		Type:              searchType,
		Start:             start,
		SearchLanguage:    lrQuery,
		InterfaceLanguage: hlQuery,
	}
}

func SearchRoute(c *gin.Context) {
	searchTerm, _ := url.QueryUnescape(c.Query("q"))
	queryParams := createSearchQueryParams(c)
	currentUrl := c.Request.URL

	if len(searchTerm) == 0 {
		c.Redirect(http.StatusPermanentRedirect, "/")
		return
	}

	if queryParams.Type == "isch" {
		searchResponse, err := ImageSearch(searchTerm, queryParams)
		if err != nil {
			log.Println(err)
		}
		imagesPageContext := createImagesPageContext(*searchResponse.ImagesPage, currentUrl)
		c.HTML(http.StatusOK, "image-search-page", imagesPageContext)
		return
	} else if queryParams.Type == "vid" {
		searchResponse, err := VideoSearch(searchTerm, queryParams)
		if err != nil {
			log.Println(err)
		}
		videosPageContext := createVideosPageContext(*searchResponse.VideosPage, currentUrl)
		c.HTML(http.StatusOK, "video-search-page", videosPageContext)

		logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal("error opening log file")
		}

		log.SetOutput(logFile)
		spew.Fdump(log.Writer(), searchResponse)
		defer logFile.Close()

		return
	} else if len(queryParams.Type) > 0 {
		// c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("https://google.com/search?q=%s&tbm=%s", searchTerm, searchType))
		return
	}

	searchResponse, err := Search(searchTerm, queryParams)

	if err != nil {
		log.Fatal(err)
	}

	if searchResponse.Type == SearchResponseError {
		log.Printf("search response error with code: %d", searchResponse.Status)
	} else if searchResponse.Type == SearchResponsePage {
		searchPageContext := createSearchPageContext(*searchResponse.SearchPage, currentUrl)
		c.HTML(http.StatusOK, "search-page", searchPageContext)
	} else if searchResponse.Type == SearchResponseCaptcha {
		captchaPageContext := createCaptchaPageContext(*searchResponse.Captcha)
		c.HTML(http.StatusOK, "captcha-page", captchaPageContext)
	}

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("error opening log file")
	}

	log.SetOutput(logFile)
	spew.Fdump(log.Writer(), searchResponse)
	defer logFile.Close()
}
