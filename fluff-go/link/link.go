package link

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Link ...
type Link struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

const linkMinLength = 2
const linkMaxLength = 32

// ValidateURL ...
func (link *Link) ValidateURL() (bool, error) {
	// TODO
	if !(strings.HasPrefix(link.URL, "http://") || strings.HasPrefix(link.URL, "https://")) {
		var sb strings.Builder
		sb.WriteString("https://")
		sb.WriteString(link.URL)
		link.URL = sb.String()
	}
	u, err := url.ParseRequestURI(link.URL)
	if err != nil {
		return false, errors.New("invalid url")
	}
	link.URL = u.String()
	return true, nil
}

// ValidateKey ...
func (link *Link) ValidateKey() (bool, error) {
	link.Key = strings.ReplaceAll(link.Key, " ", "")
	if len(link.Key) < linkMinLength {
		return false, fmt.Errorf("min length %d", linkMinLength)
	}
	if len(link.Key) > linkMaxLength {
		return false, fmt.Errorf("max length %d", linkMaxLength)
	}
	for _, c := range link.Key {
		// 0-9, a-z, A-Z
		if !((c > 47 && c < 58) || (c > 64 && c < 91) || (c > 96 && c < 123)) {
			return false, errors.New("chars")
		}
	}

	return true, nil
}
