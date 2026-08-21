package story

type MediaType int8

const (
	MediaImage MediaType = 1
	MediaAudio MediaType = 2
)

type Media struct {
	ID      int64     `json:"id"`
	StoryID int64     `json:"storyId"`
	Type    MediaType `json:"type"`
	URL     string    `json:"url"`
	Order   int       `json:"order"`
}
