package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/yuchida-tamu/git-workout-api/internal/record"
)

type PostRecordRequest struct {
	MessageBody string `json:"message_body" validate:"required"`
	Author      string `json:"author" validate:"required"`
}

func convertPostRecordRequestToRecord(r PostRecordRequest) record.Record {
	// get current time
	date := time.Now().Format("2006-01-02")
	return record.Record{
		DateCreated: date,
		MessageBody: r.MessageBody,
		Author:      r.Author,
	}
}

type RecordService interface {
	GetRecordByAuthor(context.Context, string) ([]record.Record, error)
	GetRecordById(context.Context, string) (record.Record, error)
	PostRecord(context.Context, record.Record) (record.Record, error)
	UpdateRecord(ctx context.Context, ID string, rcd record.Record) (record.Record, error)
	DeleteRecord(context.Context, string) error
}

func (h *Handler) PostRecord(w http.ResponseWriter, r *http.Request) {
	var record PostRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		return
	}

	validate := validator.New()
	err := validate.Struct(record)
	if err != nil {
		http.Error(w, "not a valid record", http.StatusBadRequest)
		return
	}

	postedRecord, err := h.Service.Record.PostRecord(r.Context(), convertPostRecordRequestToRecord(record))
	if err != nil {
		log.Print(err)
		return
	}

	if err := json.NewEncoder(w).Encode(postedRecord); err != nil {
		panic(err)
	}
}

func (h *Handler) GetRecordByAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	records, err := h.Service.Record.GetRecordByAuthor(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(records); err != nil {
		panic(err)
	}
}

func (h *Handler) GetRecordById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record, err := h.Service.Record.GetRecordById(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(record); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var record record.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		return
	}

	validate := validator.New()
	err := validate.Struct(record)
	if err != nil {
		http.Error(w, "not a valid record", http.StatusBadRequest)
		return
	}

	record, err = h.Service.Record.UpdateRecord(r.Context(), id, record)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(record); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Service.Record.DeleteRecord(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Response{message: "Successfully deleted"}); err != nil {
		panic(err)
	}
}
