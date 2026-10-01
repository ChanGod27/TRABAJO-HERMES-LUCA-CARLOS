package main

import (
	"ServidorEntregable/db"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func createUsuario(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	var newUsuario db.Usuario

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

	params := db.CreateUserParams{
		Nombre:     newUsuario.Nombre,
		Apellido:   newUsuario.Apellido,
		Email:      newUsuario.Email,
		Contrasena: newUsuario.Contrasena,
	}

	id, err := q.CreateUser(ctx, params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // El error 500
	}

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func getUserById(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {

	idStr := r.PathValue("id")

	id64, err := strconv.ParseInt(idStr, 10, 32)

	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	id := int32(id64)

	newUsuario, err := q.GetUserByID(ctx, id)

	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "aplication/json")
	json.NewEncoder(w).Encode(newUsuario)

}

func getListMateriasByUsuario(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	idStr := r.PathValue("id")

	id64, err := strconv.ParseInt(idStr, 10, 32)

	if err != nil {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	id := int32(id64)

	listaMateria, err := q.ListMateriasByUsuario(ctx, id)

	if err != nil {
		http.Error(w, "Usuario para listar materias no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "aplication/json")
	json.NewEncoder(w).Encode(listaMateria)

}

func putUpdateMateria(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {

	var UpdateMateria db.ListaTarea

	err := json.NewDecoder(r.Body).Decode(&UpdateMateria)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(UpdateMateria.Titulo) == "" ||
		UpdateMateria.IDLista <= 0 ||
		UpdateMateria.IDUsuario <= 0 {
		http.Error(w, "Todos los campos son obligatorios y no pueden contener solo espacios", http.StatusBadRequest)
		return
	}

	params := db.UpdateMateriaParams{
		IDLista: UpdateMateria.IDLista,
		Titulo:  UpdateMateria.Titulo,
	}

	// id, err := q.UpdateMateria(ctx, params) // Preguntar que hago con el id

	// Analizar posible MateriaByID
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Evaluar posible error
		return
	}

}

func main() {
	//Conexion con la base de datos
	ctx := context.Background()

	connStr := "postgres://trabajo_web:carlos-luca-hermes@localhost:5432/go_tasks_db?sslmode=disable"

	conn, errOpen := abridDB(ctx, connStr)
	if errOpen != nil {
		fmt.Printf("error crítico al conectar a la base de datos: %v", errOpen)
	}
	defer conn.Close(ctx)

	//Instanciar objeto de la base de datos (userRepository)
	queries := db.New(conn)

	//Implementacion del servidor de archivos
	staticDir := "./static"

	fileServer := http.FileServer(http.Dir(staticDir))

	http.Handle("/", fileServer)

	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}

	//Handlers que conectan las operaciones crud

	http.HandleFunc("/usuarios")
}
