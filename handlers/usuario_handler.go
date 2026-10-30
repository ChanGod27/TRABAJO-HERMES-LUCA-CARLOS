package handlers

import (
	"ServidorEntregable/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// ==========================================
// OPERACIONES DE USUARIO
// ==========================================

func createUsuario(w http.ResponseWriter, r *http.Request, q *db.Queries) {
	var newUsuario db.CreateUserParams

	err := json.NewDecoder(r.Body).Decode(&newUsuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newUsuario.Nombre) == "" ||
		strings.TrimSpace(newUsuario.Apellido) == "" ||
		strings.TrimSpace(newUsuario.Email) == "" ||
		strings.TrimSpace(newUsuario.Contrasena) == "" {
		http.Error(w, "Todos los campos son obligatorios y no pueden contener solo espacios", http.StatusBadRequest)
		return
	}

	userCreated, err := q.CreateUser(r.Context(), newUsuario)
	if err != nil {
		log.Printf("ERROR BD al crear usuario: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userCreated)
}

func getUserByEmail(w http.ResponseWriter, r *http.Request, q *db.Queries) {
	email := r.PathValue("email")

	if strings.TrimSpace(email) == "" {
		http.Error(w, "El email no puede estar vacío", http.StatusBadRequest)
		return
	}

	newUsuario, err := q.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUsuario)
}

func getUserById(w http.ResponseWriter, r *http.Request, q *db.Queries) {
	idStr := r.PathValue("id")

	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	id := int32(id64)

	newUsuario, err := q.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUsuario)
}