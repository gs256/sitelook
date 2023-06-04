package search

import "github.com/PuerkitoBio/goquery"

func selectionEmpty(selection *goquery.Selection) bool {
	return selection.Length() == 0
}

func findSingle(selection *goquery.Selection, selector string) *goquery.Selection {
	return selection.FindMatcher(goquery.Single(selector))
}

func hasInside(selection *goquery.Selection, selector string) bool {
	return !selectionEmpty(findSingle(selection, selector))
}
