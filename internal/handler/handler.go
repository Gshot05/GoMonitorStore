// internal/handler/handler.go

package handler

import (
	"encoding/json"
	"log"
	"myapp/internal/auth"
	"myapp/internal/models"
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
	var req models.UserRegisterRequest
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
	var req models.UserLoginRequest
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
	var req models.DisplayRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Diagonal <= 0 || req.Resolution == "" || req.Type == "" {
		http.Error(w, "All fields are required and must be valid", http.StatusBadRequest)
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
	var req models.MonitorRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request: unable to decode JSON", http.StatusBadRequest)
		log.Printf("Decode error: %v", err)
		return
	}

	if req.Voltage <= 0 || req.DisplayID <= 0 {
		http.Error(w, "All fields are required and must be valid", http.StatusBadRequest)
		log.Printf("Validation error: Voltage=%f, DisplayID=%d", req.Voltage, req.DisplayID)
		return
	}

	err = h.monitorService.AddMonitorWithDisplayID(req.Voltage, req.GsyncPrem, req.Curved, req.DisplayID)
	if err != nil {
		log.Printf("Error adding monitor: %v", err)
		http.Error(w, "Error adding monitor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
