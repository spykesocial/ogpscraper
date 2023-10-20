package metaparser

import (
	"net/url"

	"github.com/gocolly/colly"
)

type MetaTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GetMetaTags parses the html page at `addr` to get it's `meta` tags.
func GetMetaTags(addr string) ([]MetaTag, error) {
	tags := []MetaTag{}

	a, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	// start collecting meta tags on cache miss
	c := colly.NewCollector()
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		t := MetaTag{}
		if e.Attr("property") != "" {
			t.Name = e.Attr("property")
		} else if e.Attr("name") != "" {
			t.Name = e.Attr("name")
		} else {
			return // don't add
		}
		t.Value = e.Attr("content")
		tags = append(tags, t)
	})
	if err := c.Visit(a.String()); err != nil {
		return nil, err
	}

	return tags, err
}
