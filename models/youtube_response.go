package models

import "time"

type YoutubeDatasId struct {
	Id string `json:"videoId"`
}
type YoutubeDatasSnippetThumbnailsDefault struct {
	Url string `json:"url"`
}
type YoutubeDatasSnippetThumbnails struct {
	Default YoutubeDatasSnippetThumbnailsDefault `json:"default"`
}
type YoutubeDatasSnippet struct {
	Title       string                        `json:"title"`
	Description string                        `json:"description"`
	PublishedAt time.Time                     `json:"publishedAt"`
	Thumbnails  YoutubeDatasSnippetThumbnails `json:"thumbnails"`
}
type YoutubeDatas struct {
	Id      YoutubeDatasId      `json:"id"`
	Snippet YoutubeDatasSnippet `json:"snippet"`
}
type YoutubeResponse struct {
	Datas []YoutubeDatas `json:"items"`
}
