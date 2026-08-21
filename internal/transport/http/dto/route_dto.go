package dto

type RouteWaypointDTO struct {
	Name         string  `json:"name"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	StayDuration int     `json:"stayDuration"`
	Order        int     `json:"order"`
}

type RouteRequest struct {
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	ThemeID       *int64            `json:"themeId"`
	City          string            `json:"city"`
	StartLat      *float64          `json:"startLat"`
	StartLng      *float64          `json:"startLng"`
	EndLat        *float64          `json:"endLat"`
	EndLng        *float64          `json:"endLng"`
	TotalDistance *int              `json:"totalDistance"`
	Duration      *int              `json:"duration"`
	Difficulty    int8              `json:"difficulty"`
	CoverImage    string            `json:"coverImage"`
	Waypoints     []RouteWaypointDTO `json:"waypoints"`
	Status        int8              `json:"status"`
}

type RouteRateRequest struct {
	Score float64 `json:"score"`
}

