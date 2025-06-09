package handlers

import (
	"encoding/json"
	"encrypted-chat-relay/models"
	"fmt"
	"net/http"
	"strings"
)

type registerRequest struct {
	UserName  string `json:"username"`
	PublicKey string `json:"public_key"`
}

func (h *httpHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// Check content type
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	// Max hody size 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	// Setup decoder and deode body to struct
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req registerRequest

	if err := dec.Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if len(req.UserName) < 6 {
		http.Error(w, "Username too short, must be longer than 6 characters.", http.StatusBadRequest)
		return
	}
	
	// Check if username already exists
	var existingUser models.User
	result := h.db.Where("user_name = ?", req.UserName).First(&existingUser)
	if result.RowsAffected > 0 {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}
	
	// Create new user
	user := models.User{
		UserName: req.UserName,
		PublicKey: req.PublicKey,
	}
	
	// Save to database
	if err := h.db.Create(&user).Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to save user: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
		"id": user.ID.String(),
	})
}
