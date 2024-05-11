package scraper

import (
	"net/url"
	"time"

	"github.com/vkhobor/go-opencv/youtube"
)

const apiForQuery = "search?q="

func apiForSearchQuery(baseUrl string, search string) string {
	return baseUrl + apiForQuery + url.QueryEscape(search)
}

type Scraper struct {
	Throttle time.Duration
	Domain   string
}

func (s Scraper) Scrape(limit int, query string) []youtube.YoutubeVideo {
	singlePageVisitor := NewCollector(s)
	defer singlePageVisitor.Stop()

	singlePageVisitor.MaxDepth = 1

	allPagesVisitor := NewCollector(s)
	defer allPagesVisitor.Stop()

	youtubeIDs := make(chan youtube.YoutubeVideo, 1)
	defer close(youtubeIDs)

	allPagesVisitor.OnVideoDetailLink(func(link string) {
		singlePageVisitor.Visit(link)
	})

	singlePageVisitor.OnYoutubeUrl(func(url string) {
		id, err := youtube.NewYoutubeIDFromUrl(url)
		if err != nil {
			return
		}
		youtubeIDs <- id
	})

	urlEncoded := url.QueryEscape(query)
	allPagesVisitor.Visit(apiForSearchQuery(s.Domain, urlEncoded))

	youtubeIDsSlice := []youtube.YoutubeVideo{}
	for item := range youtubeIDs {
		if limit > 0 {
			youtubeIDsSlice = append(youtubeIDsSlice, item)
			limit--
		} else {

			break
		}
	}

	return youtubeIDsSlice
}
