package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net-http/internal/repository"
	db "net-http/internal/repository"
	"net/http"
	"strconv"
)

type APIServer struct {
	listenAddr string
	store      *db.Queries
}

func NewAPIServer(listenAddr string, store *db.Queries) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := http.NewServeMux()
	router.HandleFunc("GET /users", makeHTTPHandleFunc(s.handleGetUsers))
	router.HandleFunc("GET /users/{id}", makeHTTPHandleFunc(s.handleGetUserByID))
	router.HandleFunc("POST /users", makeHTTPHandleFunc(s.handleCreateUser))
	router.HandleFunc("PUT /users", makeHTTPHandleFunc(s.handleUpdateUser))
	router.HandleFunc("DELETE /users/{id}", makeHTTPHandleFunc(s.handleDeleteUser))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.FindAllUsers(r.Context())

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	user, err := s.store.FindUser(r.Context(), id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	createUserReq := &CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}

	defer r.Body.Close()
	user := repository.CreateUserParams{Name: createUserReq.Name, Email: createUserReq.Email}
	i, err := s.store.CreateUser(r.Context(), user)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, i)
}

func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	updateUserReq := &UpdateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(updateUserReq); err != nil {
		return err
	}

	defer r.Body.Close()
	user := repository.UpdateUsersParams{ID: updateUserReq.ID, Name: updateUserReq.Name, Email: updateUserReq.Email}
	i, err := s.store.UpdateUsers(r.Context(), user)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, i)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteUser(r.Context(), id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int64{"delete": id})
}

type CreateUserRequest struct {
	Name  string
	Email string
}

type UpdateUserRequest struct {
	ID    int64
	Name  string
	Email string
}

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
