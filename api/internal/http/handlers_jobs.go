package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/BimaPDev/SignalStack/api/internal/model"
	"github.com/BimaPDev/SignalStack/api/internal/repo"
)

type JobsHandler struct {
	Repo *repo.JobsRepo
	Log  *slog.Logger
}

func (h *JobsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.Type == "" {
		http.Error(w, `{"error":"req.Type not found"}`, http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	res, err := h.Repo.Insert(r.Context(), userID, req)
	if err != nil {
		h.Log.Error("insert failed", "err", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
func (h *JobsHandler) List(w http.ResponseWriter, r *http.Request) {
	//var req model.ListJobsResponse
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
	//	return
	//}
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	q := r.URL.Query()
	status := q.Get("status")
	limitStr := q.Get("limit")
	offsetStr := q.Get("offset")
	limit := 20
	if limitStr != "" {
		parsedStat, err := strconv.Atoi(limitStr)
		if err != nil || parsedStat <= 0 {
			http.Error(w, `{"error": "invalid limit"}`, http.StatusBadRequest)
			return
		}
		limit = parsedStat
	}
	offset := 0
	if offsetStr != "" {
		parsedOff, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOff < 0 {
			http.Error(w, `{"error": "invalid offset"}`, http.StatusBadRequest)
			return
		}
		offset = parsedOff
	}
	res, err := h.Repo.List(r.Context(), userID, status, limit, offset)
	if err != nil {
		h.Log.Error("List Fail", "err", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
func (h *JobsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	read_id := r.PathValue("id")
	res, err := h.Repo.GetByID(r.Context(), userID, read_id)
	if err != nil {
		h.Log.Error("Fetch Fail", "err", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}
