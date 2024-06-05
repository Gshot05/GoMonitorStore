// internal/models/models.go

package models

type DisplayRequest struct {
	Diagonal   float64 `json:"diagonal"`
	Resolution string  `json:"resolution"`
	Type       string  `json:"type"`
	Gsync      bool    `json:"gsync"`
}

type MonitorRequest struct {
	Voltage   float64 `json:"voltage"`
	GsyncPrem bool    `json:"gsyncPrem"`
	Curved    bool    `json:"curved"`
	DisplayID int64   `json:"displayID"`
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
