package opengraph

import "github.com/yashdiniz/ogpscraper/internal/metaparser"

const (
	UrlKey    = "og:url"
	TitleKey  = "og:title"
	DescKey   = "og:description"
	ImageKey  = "og:image"
	TypeKey   = "og:type"
	LocaleKey = "og:locale"
)

type Result struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Type        string `json:"type"`
	Locale      string `json:"locale"`
}

func GetOGPResult(tags []metaparser.MetaTag) Result {
	res := Result{}

	for _, tag := range tags {
		switch tag.Name {
		case UrlKey:
			res.Url = tag.Value
		case TitleKey:
			res.Title = tag.Value
		case DescKey:
			res.Description = tag.Value
		case ImageKey:
			res.Image = tag.Value
		case TypeKey:
			res.Type = tag.Value
		case LocaleKey:
			res.Locale = tag.Value
		}
	}

	return res
}
