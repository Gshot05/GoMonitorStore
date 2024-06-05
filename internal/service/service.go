// internal/service/service.go

package service

import (
	"database/sql"
	"log"

	"github.com/nats-io/nats.go"
)

type DisplayService struct {
	db *sql.DB
	nc *nats.Conn
}

type MonitorService struct {
	db *sql.DB
	nc *nats.Conn
}

func NewDisplayService(db *sql.DB, nc *nats.Conn) *DisplayService {
	return &DisplayService{
		db: db,
		nc: nc}
}

func NewMonitorService(db *sql.DB, nc *nats.Conn) *MonitorService {
	return &MonitorService{
		db: db,
		nc: nc}
}

func (s *DisplayService) AddDisplay(diagonal float64, resolution, displayType string, gsync bool) error {
	query := "INSERT INTO Type_Display (Name_Diagonal, Name_Resolution, Type_Type, Type_Gsync) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, diagonal, resolution, displayType, gsync)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	message := "Добавлен новый дисплей"
	err = s.nc.Publish("log", []byte(message))
	if err != nil {
		log.Printf("Error publishing message: %v", err)
		return err
	}

	return nil
}

func (s *MonitorService) AddMonitorWithDisplayID(voltage float64, gsyncPremium bool, curved bool, displayID int64) error {
	query := "INSERT INTO Type_Monitor (Name_Voltage, Name_Gsync_Prem, Name_Curved, Type_Display_ID) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, voltage, gsyncPremium, curved, displayID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	message := "Добавлен новый монитор"
	err = s.nc.Publish("log", []byte(message))
	if err != nil {
		log.Printf("Error publishing message: %v", err)
		return err
	}

	return nil
}
