package user

import "time"

type Badge struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Icon          string    `json:"icon"`
	ConditionType string    `json:"conditionType"`
	ConditionValue int      `json:"conditionValue"`
	CreatedAt     time.Time `json:"createdAt"`
}
