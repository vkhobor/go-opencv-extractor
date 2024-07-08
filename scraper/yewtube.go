package scraper

import (
	"log/slog"
	"strings"
	"sync/atomic"

	"github.com/gocolly/colly"
)

type myCollyCollector struct {
	*colly.Collector
	stopped *atomic.Bool
}

func NewCollector(s Scraper) myCollyCollector {
	defaultColl := colly.NewCollector(
		colly.AllowedDomains(s.Domain),
		// TODO colly.Async(true) enable or check if stopped needs to be atomic
		colly.UserAgent(""),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	defaultColl.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: s.Throttle,
	})

	var atomicBool atomic.Bool
	atomicBool.Store(false)

	v := myCollyCollector{defaultColl, &atomicBool}
	defaultColl.OnRequest(func(r *colly.Request) {
		if v.Stopped() {
			slog.Debug("Aborting request", "url", r.URL)
			r.Abort()
		}
		slog.Debug("Request", "url", r.URL)
	})

	defaultColl.OnResponse(func(r *colly.Response) {
		slog.Debug("Response", "url", r.Request.URL)
	})

	slog.Debug("Created new collector", "throttle", s.Throttle, "domain", s.Domain, "stopped", v.Stopped())
	return v
}

func (c myCollyCollector) Stopped() bool {
	return c.stopped.Load()
}

func (c myCollyCollector) Stop() {
	slog.Debug("Stopping collector")
	c.stopped.Store(true)
}

func (c myCollyCollector) OnVideoDetailLink(handler func(link string)) {
	c.OnHTML("div.thumbnail", func(e *colly.HTMLElement) {
		if c.Stopped() {
			return
		}

		href := e.ChildAttr("a", "href")
		if !strings.Contains(href, "watch") {
			return
		}

		handler(e.Request.AbsoluteURL(href))
	})

	c.OnHTML(`div.page-next-container`, func(e *colly.HTMLElement) {
		if c.Stopped() {
			return
		}

		link := e.ChildAttr("a", "href")
		if link == "" {
			return
		}

		e.Request.Visit(link)
	})
}

func (c myCollyCollector) OnYoutubeUrl(handler func(url string)) {
	c.OnHTML("a#link-yt-watch", func(e *colly.HTMLElement) {
		if c.Stopped() {
			return
		}

		href := e.Attr("href")
		if href == "" {
			return
		}

		handler(href)
	})
}
