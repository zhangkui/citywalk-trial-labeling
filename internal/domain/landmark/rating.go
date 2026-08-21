package landmark

type Rating struct {
	ID         int64   `json:"id"`
	LandmarkID int64   `json:"landmarkId"`
	UserID     int64   `json:"userId"`
	Score      float64 `json:"score"`
}
