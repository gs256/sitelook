package search

import (
	"net/url"
	"strconv"
	"strings"
)

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

func makeUrlTitle(itemUrl string) (string, error) {
	itemUrl, _ = url.PathUnescape(itemUrl)
	itemUrl = strings.TrimPrefix(itemUrl, "https://")
	itemUrl = strings.TrimPrefix(itemUrl, "http://")
	itemUrl = strings.TrimPrefix(itemUrl, "www.")
	itemUrl = strings.TrimSuffix(itemUrl, "/")
	return itemUrl, nil
}

func createSearchResultContext(searchResult SearchResult) SearchResultContext {
	urlTitle, _ := makeUrlTitle(searchResult.Url)

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

func createPaginationContext(pagination MultiPagePagination, currentUrl *url.URL) MultiPagePaginationContext {
	pageLinks := make([]PageLinkContext, len(pagination.PageLinks))
	query := currentUrl.Query()

	for i := 0; i < len(pagination.PageLinks); i++ {
		link := pagination.PageLinks[i]
		query.Set("start", strconv.Itoa(link.Offset))
		pageLinks[i] = createPageLinkContext(link, createHref(currentUrl, query))
	}

	query.Set("start", strconv.Itoa(pagination.PreviousOffset))
	previousUrl := createHref(currentUrl, query)
	query.Set("start", strconv.Itoa(pagination.NextOffset))
	nextUrl := createHref(currentUrl, query)

	return MultiPagePaginationContext{
		Visible:            len(pageLinks) > 0,
		Type:               PaginationTypeMultiPage,
		PageLinks:          pageLinks,
		PreviousUrl:        previousUrl,
		PreviousLinkActive: len(pageLinks) > 0 && !pageLinks[0].IsCurrent,
		NextUrl:            nextUrl,
		NextLinkActive:     len(pageLinks) > 0 && !pageLinks[len(pageLinks)-1].IsCurrent,
	}
}

func createSinglePagePaginationContext(pagination SinglePagePagination, currentUrl *url.URL) SinglePagePaginationContext {
	query := currentUrl.Query()
	query.Set("start", strconv.Itoa(pagination.PreviousOffset))
	previousLinkHref := createHref(currentUrl, query)
	query.Set("start", strconv.Itoa(pagination.NextOffset))
	nextLinkHref := createHref(currentUrl, query)
	query.Set("start", "0")
	firstPageLinkHref := createHref(currentUrl, query)

	return SinglePagePaginationContext{
		Visible:              true,
		Type:                 PaginationTypeSinglePage,
		FirstPageLinkPresent: pagination.PreviousLinkPresent && pagination.PreviousOffset != 0,
		FirstPageUrl:         firstPageLinkHref,
		PreviousLinkPresent:  pagination.PreviousLinkPresent,
		PreviousUrl:          previousLinkHref,
		NextLinkPresent:      pagination.NextLinkPresent,
		NextUrl:              nextLinkHref,
		CurrentTitle:         pagination.CurrentTitle,
	}
}

func getCurrentSearchType(queryParam string) string {
	if queryParam == "isch" {
		return SearchTypeImages
	} else if queryParam == "vid" {
		return SearchTypeVideos
	} else {
		return SearchTypeAll
	}
}

func createNavigationContext(currentUrl *url.URL) SearchNavigationContext {
	query := currentUrl.Query()
	tbm := query.Get("tbm")
	searchType := getCurrentSearchType(tbm)

	query.Del("start")
	query.Del("tbm")
	allSearchHref := createHref(currentUrl, query)
	query.Set("tbm", "isch")
	imageSearchHref := createHref(currentUrl, query)
	query.Set("tbm", "vid")
	videoSearchHref := createHref(currentUrl, query)

	return SearchNavigationContext{
		CurrentSearchType: searchType,
		SearchQueryParam:  tbm,
		AllSearchHref:     allSearchHref,
		ImageSearchHref:   imageSearchHref,
		VideoSearchHref:   videoSearchHref,
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

func createImageResultContext(imageResult ImageResult) ImageResultContext {
	return ImageResultContext{
		Title:         imageResult.Title,
		UrlTitle:      imageResult.UrlTitle,
		ImageSrc:      imageResult.ImageSrc,
		TitleLinkHref: imageResult.TitleLinkHref,
		ImageLinkHref: imageResult.ImageLinkHref,
	}
}

func createVideoResultContext(videoResult VideoResult) VideoResultContext {
	urlTitle, _ := makeUrlTitle(videoResult.TitleLinkHref)

	return VideoResultContext{
		Title:         videoResult.Title,
		UrlTitle:      urlTitle,
		ImageSrc:      videoResult.ImageSrc,
		TitleLinkHref: videoResult.TitleLinkHref,
		Description:   videoResult.Description,
	}
}

func createImagesPageContext(imagesPage ImagesPage, currentUrl *url.URL) ImagesPageContext {
	imageResults := make([]ImageResultContext, len(imagesPage.ImageResults))

	for i := 0; i < len(imagesPage.ImageResults); i++ {
		imageResults[i] = createImageResultContext(imagesPage.ImageResults[i])
	}

	return ImagesPageContext{
		SearchTerm:       imagesPage.SearchTerm,
		ImageResults:     imageResults,
		Pagination:       createSinglePagePaginationContext(imagesPage.Pagination, currentUrl),
		Navigation:       createNavigationContext(currentUrl),
		SearchCorrection: SearchCorrectionContext{},
	}
}

func createVideosPageContext(videosPage VideosPage, currentUrl *url.URL) VideosPageContext {
	videoResults := make([]VideoResultContext, len(videosPage.VideoResults))

	for i := 0; i < len(videosPage.VideoResults); i++ {
		videoResults[i] = createVideoResultContext(videosPage.VideoResults[i])
	}

	return VideosPageContext{
		SearchTerm:       videosPage.SearchTerm,
		VideoResults:     videoResults,
		Pagination:       createSinglePagePaginationContext(videosPage.Pagination, currentUrl),
		Navigation:       createNavigationContext(currentUrl),
		SearchCorrection: SearchCorrectionContext{},
	}
}

func createEmptySearchPageContext() SearchPageContext {
	return SearchPageContext{}
}

func createCaptchaPageContext(captchaPage CaptchaPage) CaptchaPageContext {
	searchUrl := getSearchUrl(captchaPage.SearchTerm, 0, "")

	return CaptchaPageContext{
		SearchRedirectUrl: searchUrl,
	}
}

func createCaptchaPage(searchTerm string) CaptchaPage {
	return CaptchaPage{
		SearchTerm: searchTerm,
	}
}

func createEmptySearchPage() SearchPage {
	return SearchPage{}
}

func createEmptyImagesPage() ImagesPage {
	return ImagesPage{
		SearchTerm:   "",
		ImageResults: []ImageResult{},
		Pagination:   SinglePagePagination{},
	}
}

func createEmptyVideosPage() VideosPage {
	return VideosPage{
		SearchTerm:   "",
		VideoResults: []VideoResult{},
		Pagination:   SinglePagePagination{},
	}
}
