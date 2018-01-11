package crawl

import (
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

// Only enqueue the root and paths matching expression.
var exp = regexp.MustCompile(`http://github\.com/kubernetes/kubernetes/tree/master(/[a-g].*)?$`)

// Extender creates the Extender implementation.
type Extender struct {
	gocrawl.DefaultExtender
}

// Visit overrides the default Visit function.
func (ext *Extender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	return nil, true
}

// Filter overrides the default Filter function.
func (ext *Extender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && exp.MatchString(ctx.NormalizedURL().String())
}

func Crawl() {
	// Set custom options.
	opts := gocrawl.NewOptions(new(Extender))
	opts.RobotUserAgent = "CCBot"
	opts.UserAgent = "TypoSpider"
	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogAll
	opts.MaxVisits = 5

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("https://github.com/kubernetes/kubernetes/tree/master")
}
