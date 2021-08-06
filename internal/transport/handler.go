package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"some-http-server/internal/types"
)

type TaskService interface {
	Create(ctx context.Context, description string) (types.QuoteResponseData, error)
	Read(ctx context.Context, id string) (types.FullQuoteData, error)
}

type TaskHandler struct {
	svc TaskService
}

func NewTaskHandler(svc TaskService) *TaskHandler {
	return &TaskHandler{
		svc: svc,
	}
}

func (t *TaskHandler) Register(r *mux.Router) {
	r.HandleFunc("/tasks", t.create).Methods(http.MethodPost)
	r.HandleFunc(fmt.Sprintf("/tasks/{id:%s}", "someRegex"), t.read).Methods(http.MethodGet)
}

type CreateTasksRequest struct {
	Description string `json:"description"`
}

type CreateTasksResponse struct {
	Task types.Task `json:"task"`
}

func (t *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, "invalid request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	task, err := t.svc.Create(r.Context(), req.Description)
	if err != nil {
		renderErrorResponse(w, "create failed", http.StatusInternalServerError)
		return
	}

	renderResponse(w,
		&CreateTasksResponse{
			Task: types.Task{
				ID:          task.ID,
				Description: task.Description,
			},
		},
		http.StatusCreated)
}

type GetTasksResponse struct {
	Task types.Task `json:"task"`
}

func (t *TaskHandler) read(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		// TODO handle
	}

	task, err := t.svc.Read(r.Context(), id)
	if err != nil {
		renderErrorResponse(w, "find failed", http.StatusNotFound)
		return
	}

	renderResponse(w,
		&GetTasksResponse{
			Task: types.Task{
				ID:          task.ID,
				Description: task.Description,
			},
		},
		http.StatusOK)
}

// UpdateTasksRequest defines the request used for updating a task.
type UpdateTasksRequest struct {
	Description string `json:"description"`
}

func (t *TaskHandler) update(w http.ResponseWriter, r *http.Request) {
	var req UpdateTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, "invalid request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	id, ok := mux.Vars(r)["id"]
	if !ok {
		// TODO handle
	}

	err := t.svc.Update(r.Context(), id, req.Description)
	if err != nil {
		renderErrorResponse(w, "update failed", http.StatusInternalServerError)
		return
	}

	renderResponse(w, &struct{}{}, http.StatusOK)
}

type DeleteTasksRequest struct {
	Id string `json:"id"`
}

func (t *TaskHandler) delete(w http.ResponseWriter, r *http.Request) {
	var req DeleteTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderErrorResponse(w, "invalid request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err := t.svc.Delete(r.Context(), req.Id)
	if err != nil {
		renderErrorResponse(w, "create failed", http.StatusInternalServerError)
		return
	}

	renderResponse(w, &struct{}{}, http.StatusCreated)
}
