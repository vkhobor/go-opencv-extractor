package scraper

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/vkhobor/go-opencv/youtube"
)

func apiForSearchQuery(host string, search string) (string, error) {
	params := url.Values{}
	params.Add("q", search)

	u := url.URL{
		Scheme:   "https",
		Host:     host,
		RawQuery: params.Encode(),
		Path:     "search",
	}
	return u.String(), nil
}

type Scraper struct {
	Throttle time.Duration
	Domain   string
}

func (s Scraper) Scrape(ctx context.Context, query string) (<-chan youtube.YoutubeVideo, error) {
	singlePageVisitor := NewCollector(s)

	singlePageVisitor.MaxDepth = 1

	allPagesVisitor := NewCollector(s)

	youtubeIDsChan := make(chan youtube.YoutubeVideo)

	allPagesVisitor.OnVideoDetailLink(func(link string) {
		singlePageVisitor.Visit(link)
	})

	singlePageVisitor.OnYoutubeUrl(func(url string) {
		id, err := youtube.NewYoutubeIDFromUrl(url)
		if err != nil {
			return
		}

		select {
		case <-ctx.Done():
			allPagesVisitor.Stop()
			singlePageVisitor.Stop()
			close(youtubeIDsChan)
		default:
			youtubeIDsChan <- id
		}
	})

	urlEncoded := url.QueryEscape(query)
	urlToScrape, err := apiForSearchQuery(s.Domain, urlEncoded)
	if err != nil {
		return nil, err
	}

	slog.Debug("Start scraping", "url", urlToScrape)
	go allPagesVisitor.Visit(urlToScrape)

	return youtubeIDsChan, nil
}
