package entity

import (
	"fmt"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/exres"
	"strings"
)

type WordpressPost struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Status        string `json:"status"`
	FeaturedMedia int    `json:"featured_media"`
}

type WordpressMedia struct {
	ID        int    `json:"id"`
	SourceURL string `json:"source_url"`
	MediaType string `json:"media_type"`
}

func NewWordpressPost(instaDetail InstagramPost, wpMedia []*exres.UploadMediaResponse) WordpressPost {
	wordpressPosts := WordpressPost{}
	wordpressPosts.Title = instaDetail.Title()
	wordpressPosts.FeaturedMedia = wpMedia[0].ID
	if instaDetail.MediaType == "IMAGE" {
		wordpressPosts.Content = fmt.Sprintf("%s%s", getImageHtml(wpMedia[0].SourceUrl), getContentHtml(instaDetail.Caption))
	} else if instaDetail.MediaType == "VIDEO" {
		wordpressPosts.Content = fmt.Sprintf("%s%s", getVideoHtml(wpMedia[0].SourceUrl), getContentHtml(instaDetail.Caption))
	} else {
		wordpressPosts.Content = getCarousel(instaDetail, wpMedia)
	}
	wordpressPosts.Status = "publish"
	return wordpressPosts
}

func getCarousel(instaDetail InstagramPost, wpMedia []*exres.UploadMediaResponse) string {
	sb := strings.Builder{}
	sb.WriteString("<div class='a-root-wordpress-instagram-slider'>")
	for _, media := range wpMedia {
		if media.MimeType == "'video/mp4" {
			sb.WriteString(getVideoHtml(media.SourceUrl))
		} else {
			sb.WriteString(getImageHtml(media.SourceUrl))
		}
	}
	sb.WriteString("</div>")
	sb.WriteString(getContentHtml(instaDetail.Caption))
	return sb.String()
}

func getVideoHtml(url string) string {
	return fmt.Sprintf(`<div><img src='%s' style='margin: 0 auto;' width='500px' height='500px'/></div>`, url)
}

func getImageHtml(url string) string {
	return fmt.Sprintf(`<div><video src='%s' style='margin: 0 auto;' width='500px' height='500px' controls>Sorry, your browser does not support embedded videos.</video></div>`, url)
}

func getContentHtml(caption string) string {
	sb := strings.Builder{}
	sb.WriteString("<p>")
	for i, row := range strings.Split(caption, "/n") {
		if i != 0 {
			sb.WriteString("<br>")
		}
		sb.WriteString(row)
	}
	sb.WriteString("</p>")
	return sb.String()
}
