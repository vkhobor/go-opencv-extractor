package scraper

import (
	"iter"
	"log/slog"
	"net/url"
	"time"

	"github.com/vkhobor/go-opencv/youtube"
)

func apiForSearchQuery(host string, search string) string {
	params := url.Values{}
	params.Add("q", search)

	u := url.URL{
		Scheme:   "https",
		Host:     host,
		RawQuery: params.Encode(),
		Path:     "search",
	}
	return u.String()
}

type Scraper struct {
	Throttle time.Duration
	Domain   string
	visitors []*myCollyCollector
}

func (s *Scraper) Stop() {
	for _, v := range s.visitors {
		v.Stop()
	}
}

type Result struct {
	youtube.YoutubeVideo
	Error error
}

func (s *Scraper) AllForQuery(query string) iter.Seq[Result] {
	return func(yield func(Result) bool) {
		s.Scrape(query, func(yv youtube.YoutubeVideo, err error, close func()) {
			if !yield(Result{
				YoutubeVideo: yv,
				Error:        err,
			}) {
				close()
			}
		})
	}
}

func (s *Scraper) Scrape(query string, onFound func(youtube.YoutubeVideo, error, func())) {
	singlePageVisitor := NewCollector(*s)
	singlePageVisitor.MaxDepth = 1

	allPagesVisitor := NewCollector(*s)

	s.visitors = append(s.visitors, &singlePageVisitor, &allPagesVisitor)

	allPagesVisitor.OnVideoDetailLink(func(link string) {
		singlePageVisitor.Visit(link)
	})

	singlePageVisitor.OnYoutubeUrl(func(url string) {
		id, err := youtube.NewYoutubeIDFromUrl(url)
		onFound(youtube.YoutubeVideo(id), err, s.Stop)
	})

	urlEncoded := url.QueryEscape(query)
	urlToScrape := apiForSearchQuery(s.Domain, urlEncoded)

	slog.Debug("Start scraping", "url", urlToScrape)
	allPagesVisitor.Visit(urlToScrape)

	allPagesVisitor.Wait()
	singlePageVisitor.Wait()
}
