package handler

import (
	"database/sql"
	"net/http"
)

type PingHandler struct {
	db *sql.DB
}

func NewPingHandler(db *sql.DB) *PingHandler {
	return &PingHandler{db: db}
}

func (handler *PingHandler) PingHandle(w http.ResponseWriter, r *http.Request) {
	err := handler.db.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
