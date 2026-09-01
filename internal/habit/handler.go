package habit

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		habits, err := h.service.List(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch habits"})
			return
		}
		json.NewEncoder(w).Encode(habits)

	case http.MethodPost:
		var input Habit
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body"})
			return
		}
		created, err := h.service.Create(ctx, input)
		if err != nil {
			status := http.StatusInternalServerError
			if err == ErrValidation {
				status = http.StatusUnprocessableEntity
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(created)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	}
}

func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid habit id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		result, err := h.service.Get(ctx, id)
		if err != nil {
			status := http.StatusInternalServerError
			if err == ErrNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(result)

	case http.MethodPut:
		var input Habit
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON body"})
			return
		}
		updated, err := h.service.Update(ctx, id, input)
		if err != nil {
			status := http.StatusInternalServerError
			switch err {
			case ErrValidation:
				status = http.StatusUnprocessableEntity
			case ErrNotFound:
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(updated)

	case http.MethodDelete:
		if err := h.service.Delete(ctx, id); err != nil {
			status := http.StatusInternalServerError
			if err == ErrNotFound {
				status = http.StatusNotFound
			}
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	}
}
