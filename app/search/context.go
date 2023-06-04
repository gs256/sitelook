package search

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

type ImagesPageContext struct {
	SearchTerm       string
	ImageResults     []ImageResultContext
	Pagination       PaginationContext
	Navigation       SearchNavigationContext
	SearchCorrection SearchCorrectionContext
}
