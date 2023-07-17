package search

import (
	"errors"
	"log"

	"github.com/PuerkitoBio/goquery"
)

type VideoResult struct {
	Title         string
	UrlTitle      string
	ImageSrc      string
	TitleLinkHref string
	Description   string
}

type VideosPage struct {
	SearchTerm   string
	VideoResults []VideoResult
	Pagination   SinglePagePagination
}

func parseVideosPage(document *goquery.Document) (VideosPage, error) {
	// TODO: extract in a separate function
	searchInput := findSingle(document.Selection, "input[name=\"q\"]")
	if selectionEmpty(searchInput) {
		return VideosPage{}, errors.New("search input not found")
	}
	searchTerm := searchInput.AttrOr("value", "")

	videoResults := make([]VideoResult, 0)

	document.Find("#main > div > div").Each(func(i int, item *goquery.Selection) {
		h3 := findSingle(item, "h3")

		// Search filters div or `Related searches` div
		if selectionEmpty(h3) {
			return
		}

		title := h3.Text()

		a := findSingle(item, "a")
		videoHref := hrefFromQuery(a.AttrOr("href", "#"))

		img := findSingle(item, "img")
		imgSrc := img.AttrOr("src", "")

		descriptionDiv := findSingle(item, "div>span").Parent()
		descriptionHtml := descriptionDiv.Text()

		videoResults = append(videoResults, VideoResult{
			Title:         title,
			UrlTitle:      videoHref, // FIXME: prettify
			ImageSrc:      imgSrc,
			TitleLinkHref: videoHref,
			Description:   descriptionHtml,
		})
	})

	if len(videoResults) == 0 {
		return VideosPage{
			SearchTerm:   searchTerm,
			VideoResults: videoResults,
			Pagination:   SinglePagePagination{},
		}, errors.New("page has no images or an error occured while parsing images")
	}

	pagination, err := parseVideoPagePagination(document)
	if err != nil {
		log.Println(err)
	}

	return VideosPage{
		SearchTerm:   searchTerm,
		VideoResults: videoResults,
		Pagination:   pagination,
	}, nil
}
