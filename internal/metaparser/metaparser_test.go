package metaparser_test

import (
	"testing"

	"github.com/yashdiniz/ogpscraper/internal/metaparser"
)

func TestGetMetaTags(t *testing.T) {
	url := "https://www.youtube.com/watch?v=0G7Wu4DnDaw"
	tags, err := metaparser.GetMetaTags(url)
	if err != nil {
		t.Error(err)
	}
	t.Log("Number of meta tags: ", len(tags))
	for _, tag := range tags {
		t.Log(tag)
	}
}
