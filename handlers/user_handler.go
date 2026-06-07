package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/okoraretega/doc_stream_server/model"
	"github.com/okoraretega/doc_stream_server/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var u model.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	u, err = h.userService.CreateUser(ctx, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		http.Error(w, "Unable to get users", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqUrl := strings.TrimPrefix(r.URL.Path, "/users/")

	if reqUrl == "" {
		http.Error(w, "Please provide an ID", http.StatusBadRequest)
		return
	}

	url, err := uuid.Parse(reqUrl)
	if err != nil {
		http.Error(w, "Provide a valid ID", http.StatusBadRequest)
		return
	}

	user, bool := h.userService.GetUserById(ctx, url)
	if bool == false {

		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqUrl := strings.TrimPrefix(r.URL.Path, "/users/")

	url, err := uuid.Parse(reqUrl)
	if err != nil {
		http.Error(w, "Please provide a valid ID", http.StatusBadRequest)
	}

	bool, err := h.userService.DeleteUser(ctx, url)
	if err != nil {
		http.Error(w, "Unknown Error occured", http.StatusInternalServerError)
	}

	if bool == false {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetUserById(w, r)
	case http.MethodDelete:
		h.DeleteUser(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
