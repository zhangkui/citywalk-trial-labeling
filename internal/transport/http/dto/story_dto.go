package dto

type StoryMediaDTO struct {
	Type int8   `json:"type"`
	URL  string `json:"url"`
	Order int   `json:"order"`
}

type StoryRequest struct {
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	CoverImage  string         `json:"coverImage"`
	LandmarkIDs []int64        `json:"landmarkIds"`
	RouteID     *int64         `json:"routeId"`
	Media       []StoryMediaDTO `json:"media"`
	Status      int8           `json:"status"`
}

type StoryLikeRequest struct {
	Add bool `json:"add"`
}

