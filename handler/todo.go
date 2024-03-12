package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		req := &model.ReadTODORequest{}

		q := r.URL.Query()
		_prevID := q.Get("prev_id")
		_size := q.Get("size")

		if _prevID == "" {
			req.PrevID = 0
		} else {
			prevID, err := strconv.ParseInt(_prevID, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
			req.PrevID = prevID
		}

		if _size == "" {
			req.Size = 5
		} else {
			size, err := strconv.ParseInt(_size, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
			req.Size = size
		}

		res, err := h.Read(r.Context(), req)
		if err != nil {
			log.Println(err)
			return
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			return
		}
	} else if r.Method == http.MethodPost {
		// r.Method == "Post" ではなく r.Method == "POST" にしないといけない
		req := &model.CreateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Println(err)
			return
		}

		if req.Subject == "" {
			w.WriteHeader(400)
			return
		}

		res, err := h.Create(r.Context(), req)
		if err != nil {
			log.Println(err)
			return
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			return
		}
	} else if r.Method == http.MethodPut {
		req := &model.UpdateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Println(err)
			return
		}

		if req.ID == 0 || req.Subject == "" {
			w.WriteHeader(400)
			return
		}

		res, err := h.Update(r.Context(), req)
		if err != nil {
			log.Println(err)
			if err.Error() == "Error: Not Found" {
				w.WriteHeader(404)
			}
			return
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			return
		}
	} else if r.Method == http.MethodDelete {
		req := &model.DeleteTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Println(err)
			return
		}

		if len(req.IDs) == 0 {
			w.WriteHeader(400)
			return
		}

		res, err := h.Delete(r.Context(), req)
		if err != nil {
			log.Println(err)
			if reflect.TypeOf(err) == reflect.TypeOf(&model.ErrNotFound{}) {
				w.WriteHeader(404)
			}
			return
		}

		if err = json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			return
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: todo}, err
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	return &model.ReadTODOResponse{TODOs: todos}, err
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: todo}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	return &model.DeleteTODOResponse{}, err
}
