package handlers

import (
	"cafeteria/internal/handlers/middleware"
	"cafeteria/internal/helpers"
	"cafeteria/internal/models"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type UserService interface {
	Register(ctx context.Context, user *models.User) (string, error)
	GetToken(ctx context.Context, username, pass string) (string, error)
}

type UserHandler struct {
	Service UserService
	Logger  *slog.Logger
}

func NewUserHandler(service UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		Service: service,
		Logger:  logger,
	}
}

func (h *UserHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /add_user", middleware.Middleware(h.AddUser))
	mux.HandleFunc("POST /add_user/", middleware.Middleware(h.AddUser))

	mux.HandleFunc("POST /get-token", middleware.Middleware(h.GetToken))
	mux.HandleFunc("POST /get-token/", middleware.Middleware(h.GetToken))

	mux.HandleFunc("GET /login", middleware.Middleware(h.Login))
	mux.HandleFunc("GET /login/", middleware.Middleware(h.Login))

	mux.HandleFunc("GET /register", middleware.Middleware(h.Register))
	mux.HandleFunc("Get /register/", middleware.Middleware(h.Register))
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request
	if err := r.ParseForm(); err != nil {
		h.Logger.Error(fmt.Sprintf("error parsing form data: %v", err))
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Extract form values
	username := r.FormValue("username")
	password := r.FormValue("password")
	password2 := r.FormValue("password2") // Add second password field
	ageStr := r.FormValue("age")
	allergens := r.FormValue("allergens") // Extract allergens field

	// Validate required fields
	if username == "" || password == "" || password2 == "" {
		h.Logger.Error("username, password, or repeat password not provided")
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("username, password, and repeat password are required"))
		return
	}

	// Check if passwords match
	if password != password2 {
		h.Logger.Error("passwords do not match")
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("passwords do not match"))
		return
	}

	// Password length validation
	if len(password) < 8 {
		h.Logger.Error("password too short")
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("password must be at least 8 characters long"))
		return
	}

	// Convert age to int
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("invalid age: %v", err))
		helpers.WriteError(w, http.StatusBadRequest, fmt.Errorf("age must be a valid integer"))
		return
	}

	// Create a User object
	isAdmin := false
	sexStr := r.FormValue("sex")
	allergensList := strings.Split(allergens, ",") // Convert the comma-separated allergens into a list

	// Convert sex to []uint8 (if needed)
	sex := []uint8(sexStr)

	user := &models.User{
		Username:  username,
		Password:  password, // Use the provided password
		Age:       age,
		IsAdmin:   isAdmin,
		Sex:       sex,
		Allergens: allergensList, // Store allergens as a list
	}

	// Register the user using the service
	token, err := h.Service.Register(r.Context(), user)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("error registering new user: %v", err))
		helpers.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Return a success response with the token
	helpers.WriteJSON(w, http.StatusOK, models.Reponse{
		Messege: "successfully registered and fetched token",
		Value:   token,
	})
}

func (h *UserHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	pass := r.FormValue("password")

	if len(username) == 0 {
		helpers.WriteError(w, http.StatusForbidden, fmt.Errorf("usesrname wasn't provided"))
		return
	}

	token, err := h.Service.GetToken(r.Context(), username, pass)
	if err != nil {
		h.Logger.Error(err.Error())
		helpers.WriteError(w, http.StatusForbidden, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, models.Reponse{Messege: "token was fetched", Value: token})

	http.SetCookie(w, &http.Cookie{
		Name:     "jwtToken",
		Value:    token,
		HttpOnly: true,
		Secure:   true,                    // Only send over HTTPS
		SameSite: http.SameSiteStrictMode, // Prevent CSRF
		Path:     "/",
		MaxAge:   86400, // Expires in 1 day
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "templates/login.html")

	t, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	t.Execute(w, body)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	dir, _ := os.Getwd()
	path := filepath.Join(dir, "templates/register.html")

	t, err := template.ParseFiles(path)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("error parsing template: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := os.ReadFile(path)
	if err != nil {
		h.Logger.Error(fmt.Sprintf("error reading template file: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, body); err != nil {
		h.Logger.Error(fmt.Sprintf("error executing template: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
