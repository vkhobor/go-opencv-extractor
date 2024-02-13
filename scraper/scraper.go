package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type MyCollyCollector struct {
	*colly.Collector
	stop *bool
}

func NewCollector() MyCollyCollector {
	defaultColl := colly.NewCollector(
		colly.AllowedDomains("yewtu.be"),
		colly.UserAgent(""),
		// colly.Debugger(&debug.LogDebugger{}),
	)

	defaultColl.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: time.Second * 5,
	})

	stop := false
	defaultColl.OnRequest(func(r *colly.Request) {
		if stop {
			r.Abort()
		}
	})

	v := MyCollyCollector{defaultColl, &stop}
	return v
}

func (c MyCollyCollector) OnVideoDetailLink(handler func(link string)) {
	c.OnHTML("div.thumbnail", func(e *colly.HTMLElement) {
		if *c.stop {
			return
		}

		href := e.ChildAttr("a", "href")
		if !strings.Contains(href, "watch") {
			return
		}

		handler(e.Request.AbsoluteURL(href))
	})

	c.OnHTML(`div.page-next-container`, func(e *colly.HTMLElement) {
		if *c.stop {
			return
		}

		link := e.ChildAttr("a", "href")
		if link == "" {
			return
		}

		e.Request.Visit(link)
	})
}

func (c MyCollyCollector) OnYoutubeUrl(handler func(url string)) {
	c.OnHTML("a#link-yt-watch", func(e *colly.HTMLElement) {
		if *c.stop {
			return
		}

		href := e.Attr("href")
		if href == "" {
			return
		}

		handler(href)
	})
}

func Scrape(search string, limit int, offset int, onUrlFound func(url string)) {
	singlePageVisitor := NewCollector()
	singlePageVisitor.MaxDepth = 1

	allPagesVisitor := NewCollector()
	page := 5 * 4
	allPagesVisitor.MaxDepth = offset/page + limit/page + 1

	videosFound := 0
	allPagesVisitor.OnVideoDetailLink(func(link string) {
		videosFound++
		if videosFound <= offset {
			return
		}

		singlePageVisitor.Visit(link)
	})

	found := 0
	singlePageVisitor.OnYoutubeUrl(func(url string) {
		onUrlFound(url)
		found++

		if found >= limit {
			*singlePageVisitor.stop = true
			*allPagesVisitor.stop = true
		}
	})

	urlEncoded := url.QueryEscape(search)
	allPagesVisitor.Visit("https://yewtu.be/search?q=" + urlEncoded)
}

func ScrapeToChannel(search string, limit int, offset int) <-chan string {
	resultUrl := make(chan string)

	go Scrape(search, limit, offset, func(url string) {
		resultUrl <- url
		fmt.Printf("Found %v\n", url)
	})
	return resultUrl
}
