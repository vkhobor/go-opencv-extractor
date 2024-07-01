package scraper

import (
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
	visitors []*myCollyCollector
}

func (s Scraper) Stop() {
	for _, v := range s.visitors {
		v.Stop()
	}
}

func (s Scraper) Scrape(query string, onFound func(youtube.YoutubeVideo, error, func())) error {
	singlePageVisitor := NewCollector(s)
	singlePageVisitor.MaxDepth = 1

	allPagesVisitor := NewCollector(s)

	s.visitors = append(s.visitors, &singlePageVisitor, &allPagesVisitor)

	allPagesVisitor.OnVideoDetailLink(func(link string) {
		singlePageVisitor.Visit(link)
	})

	singlePageVisitor.OnYoutubeUrl(func(url string) {
		id, err := youtube.NewYoutubeIDFromUrl(url)
		onFound(youtube.YoutubeVideo(id), err, s.Stop)
	})

	urlEncoded := url.QueryEscape(query)
	urlToScrape, err := apiForSearchQuery(s.Domain, urlEncoded)
	if err != nil {
		return err
	}

	slog.Debug("Start scraping", "url", urlToScrape)
	allPagesVisitor.Visit(urlToScrape)

	return nil
}
