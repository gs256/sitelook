package search

const (
	SearchTypeAll    = "All"
	SearchTypeImages = "Images"
	SearchTypeVideos = "Videos"
)

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

type MultiPagePaginationContext struct {
	Visible            bool
	Type               string
	PageLinks          []PageLinkContext
	PreviousUrl        string
	PreviousLinkActive bool
	NextUrl            string
	NextLinkActive     bool
}

type SearchNavigationContext struct {
	CurrentSearchType string
	SearchQueryParam  string
	AllSearchHref     string
	ImageSearchHref   string
	VideoSearchHref   string
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
	Pagination       SinglePagePaginationContext
	Navigation       SearchNavigationContext
	SearchCorrection SearchCorrectionContext
}

type CaptchaPageContext struct {
	SearchRedirectUrl string
}

type ImageResultContext struct {
	Title         string
	UrlTitle      string
	ImageSrc      string
	TitleLinkHref string
	ImageLinkHref string
}

type VideoResultContext struct {
	Title         string
	UrlTitle      string
	ImageSrc      string
	TitleLinkHref string
	Description   string
}

type SinglePagePaginationContext struct {
	Visible              bool
	Type                 string
	FirstPageLinkPresent bool
	FirstPageUrl         string
	PreviousLinkPresent  bool
	PreviousUrl          string
	NextLinkPresent      bool
	NextUrl              string
	CurrentTitle         string
}

type ImagesPageContext struct {
	SearchTerm       string
	ImageResults     []ImageResultContext
	Pagination       SinglePagePaginationContext
	Navigation       SearchNavigationContext
	SearchCorrection SearchCorrectionContext
}

type VideosPageContext struct {
	SearchTerm       string
	VideoResults     []VideoResultContext
	Pagination       SinglePagePaginationContext
	Navigation       SearchNavigationContext
	SearchCorrection SearchCorrectionContext
}
