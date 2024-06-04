package handler

import (
	"encoding/json"
	"myapp/internal/auth"
	"myapp/internal/service"
	"net/http"
)

type Handler struct {
	authService    *auth.AuthService
	displayService *service.DisplayService
	monitorService *service.MonitorService
}

func NewHandler(
	authService *auth.AuthService,
	displayService *service.DisplayService,
	monitorService *service.MonitorService) *Handler {

	return &Handler{
		authService:    authService,
		displayService: displayService,
		monitorService: monitorService,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = h.authService.Register(req.Username, req.Password, req.Email)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) AddDisplay(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Diagonal   float64 `json:"diagonal"`
		Resolution string  `json:"resolution"`
		Type       string  `json:"type"`
		Gsync      bool    `json:"gsync"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = h.displayService.AddDisplay(req.Diagonal, req.Resolution, req.Type, req.Gsync)
	if err != nil {
		http.Error(w, "Error adding display", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) AddMonitor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Voltage   float64 `json:"voltage"`
		GsyncPrem bool    `json:"gsync_prem"`
		Curved    bool    `json:"curved"`
		DisplayID int64   `json:"display_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request: unable to decode JSON", http.StatusBadRequest)
		return
	}

	err = h.monitorService.AddMonitorWithDisplayID(req.Voltage, req.GsyncPrem, req.Curved, req.DisplayID)
	if err != nil {
		http.Error(w, "Error adding monitor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
