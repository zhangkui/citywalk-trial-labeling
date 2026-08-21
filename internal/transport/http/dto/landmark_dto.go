package dto

type LandmarkRequest struct {
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	CategoryID  *int64  `json:"categoryId"`
	Description string  `json:"description"`
	CoverImage  string  `json:"coverImage"`
}

type LandmarkRateRequest struct {
	Score float64 `json:"score"`
}

