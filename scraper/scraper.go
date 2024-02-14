package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/vkhobor/go-opencv/importing"
)

type MyCollyCollector struct {
	*colly.Collector
	stopped *bool
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
		if *c.stopped {
			return
		}

		href := e.ChildAttr("a", "href")
		if !strings.Contains(href, "watch") {
			return
		}

		handler(e.Request.AbsoluteURL(href))
	})

	c.OnHTML(`div.page-next-container`, func(e *colly.HTMLElement) {
		if *c.stopped {
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
		if *c.stopped {
			return
		}

		href := e.Attr("href")
		if href == "" {
			return
		}

		handler(href)
	})
}

func Scrape(search string, limit int, offset int, onYoutubeIdFound func(url string)) {
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
		id, err := importing.YoutubeParser(url)
		if err != nil {
			return
		}

		onYoutubeIdFound(id)
		found++

		if found >= limit {
			*singlePageVisitor.stopped = true
			*allPagesVisitor.stopped = true
		}
	})

	urlEncoded := url.QueryEscape(search)
	allPagesVisitor.Visit("https://yewtu.be/search?q=" + urlEncoded)
}

func ScrapeToChannel(search string, limit int, offset int) <-chan string {
	resultUrl := make(chan string)

	if limit == 0 {
		close(resultUrl)
		return resultUrl
	}

	go func() {
		Scrape(search, limit, offset, func(url string) {
			resultUrl <- url
			fmt.Printf("Found %v\n", url)
		})

		close(resultUrl)
	}()
	return resultUrl
}
