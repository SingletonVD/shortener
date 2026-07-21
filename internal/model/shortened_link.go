package model

import "net/url"

type ShortenedLink struct {
	Short    string
	FullLink url.URL
}
