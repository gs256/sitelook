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
	query.Set("start", strconv.Itoa(pagination.NextOffset))
	nextUrl := createHref(currentUrl, query)

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

func createImageResultContext(imageResult ImageResult) ImageResultContext {
	return ImageResultContext{
		Title:         imageResult.Title,
		UrlTitle:      imageResult.UrlTitle,
		ImageSrc:      imageResult.ImageSrc,
		TitleLinkHref: imageResult.TitleLinkHref,
		ImageLinkHref: imageResult.ImageLinkHref,
	}
}

func createImagesPageContext(imagesPage ImagesPage, currentUrl *url.URL) ImagesPageContext {
	imageResults := make([]ImageResultContext, len(imagesPage.ImageResults))

	for i := 0; i < len(imagesPage.ImageResults); i++ {
		imageResults[i] = createImageResultContext(imagesPage.ImageResults[i])
	}

	return ImagesPageContext{
		SearchTerm:       "FIXME",
		ImageResults:     imageResults,
		Pagination:       PaginationContext{},
		Navigation:       SearchNavigationContext{},
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
	}
}
