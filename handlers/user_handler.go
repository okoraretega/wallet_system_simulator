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
	var u model.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	u, err = h.userService.CreateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.userService.GetAllUsers()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	reqUrl := strings.TrimPrefix(r.URL.Path, "/users/")

	url, err := uuid.Parse(reqUrl)
	if err != nil {
		http.Error(w, "Provide a valid ID", http.StatusBadRequest)
	}

	user, bool := h.userService.GetUserById(url)
	if bool == false {

		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	reqUrl := strings.TrimPrefix(r.URL.Path, "/users/")

	url, err := uuid.Parse(reqUrl)
	if err != nil {
		http.Error(w, "Please provide a valid ID", http.StatusBadRequest)
	}

	bool, err := h.userService.DeleteUser(url)
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
