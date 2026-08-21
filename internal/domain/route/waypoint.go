package route

type Waypoint struct {
	ID           int64   `json:"id"`
	RouteID      int64   `json:"routeId"`
	Name         string  `json:"name"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	StayDuration int     `json:"stayDuration"`
	Order        int     `json:"order"`
}
