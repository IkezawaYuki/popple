package entity

import (
	"net/url"
	"strings"
)

//type InstagramPost struct {
//	ID              string
//	Caption         string
//	MediaType       string
//	MediaURL        string
//	Timestamp       time.Time
//	ChildrenID      []string
//	ChildrenContent []ChildMedia
//}

type ChildMedia struct {
	ID        string
	MediaURL  string
	MediaType string
}

func (i *InstagramPost) FileName() (string, error) {
	parsedURL, err := url.Parse(i.MediaURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Path, nil
}

func (i *InstagramPost) Title() string {
	return strings.Split(i.Caption, " ")[0]
}

func (i *ChildMedia) FileName() (string, error) {
	parsedURL, err := url.Parse(i.MediaURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Path, nil
}

type InstagramPosts struct {
	ID    string `json:"id"`
	Media struct {
		Data []InstagramPost `json:"data"`
	} `json:"media"`
}

type InstagramPost struct {
	ID        string                `json:"id"`
	Permalink string                `json:"permalink"`
	Caption   string                `json:"caption,omitempty"`
	Timestamp string                `json:"timestamp"`
	MediaType string                `json:"media_type"`
	MediaURL  string                `json:"media_url"`
	Children  InstagramPostChildren `json:"children"`
}

type InstagramPostChildren struct {
	Data []struct {
		MediaType string `json:"media_type"`
		MediaURL  string `json:"media_url"`
		ID        string `json:"id"`
	} `json:"data"`
}
