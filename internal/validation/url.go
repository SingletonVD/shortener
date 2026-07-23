package validation

import "net/url"

func ValidateRawLink(rawLink string) bool {
	url, err := url.Parse(rawLink)

	if err != nil {
		return false
	}

	linkValid := !(url.Host == "" || (url.Scheme != "http" && url.Scheme != "https"))

	return linkValid
}
