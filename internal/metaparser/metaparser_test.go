package metaparser_test

import (
	"testing"

	"github.com/yashdiniz/ogpscraper/internal/metaparser"
)

func TestGetMetaTags(t *testing.T) {
	url := "https://www.livemint.com/news/world/israelhamas-war-day-16-gaza-sees-most-violent-night-yet-as-bombardment-kills-400-people-in-24-hours-china-biden-11698028685643.html"
	tags, err := metaparser.GetMetaTags(url)
	if err != nil {
		t.Error(err)
	}
	t.Log("Number of meta tags: ", len(tags))
	for _, tag := range tags {
		t.Log(tag)
	}
}
