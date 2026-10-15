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

// ==========================================
// 1. OPERACIONES DE USUARIO
// ==========================================

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func getUserByEmail(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	email := r.PathValue("email")

	if strings.TrimSpace(email) == "" {
		http.Error(w, "El email no puede estar vacío", http.StatusBadRequest)
		return
	}

	newUsuario, err := q.GetUserByEmail(ctx, email)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUsuario)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUsuario)
}

// ==========================================
// 2. OPERACIONES DE MATERIA (db.ListaTarea)
// ==========================================

func createMateria(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	var newMateria db.ListaTarea

	err := json.NewDecoder(r.Body).Decode(&newMateria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newMateria.Titulo) == "" || newMateria.IDUsuario <= 0 {
		http.Error(w, "Todos los campos son obligatorios y no pueden contener solo espacios", http.StatusBadRequest)
		return
	}

	params := db.CreateMateriaParams{
		IDUsuario: newMateria.IDUsuario,
		Titulo:    newMateria.Titulo,
	}

	materiaCreada, err := q.CreateMateria(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(materiaCreada)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listaMateria)
}

func putUpdateMateria(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	var MateriaRecibida db.ListaTarea

	err := json.NewDecoder(r.Body).Decode(&MateriaRecibida)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(MateriaRecibida.Titulo) == "" ||
		MateriaRecibida.IDLista <= 0 ||
		MateriaRecibida.IDUsuario <= 0 {
		http.Error(w, "Todos los campos son obligatorios y no pueden contener solo espacios", http.StatusBadRequest)
		return
	}

	params := db.UpdateMateriaParams{
		IDLista: MateriaRecibida.IDLista,
		Titulo:  MateriaRecibida.Titulo,
	}

	MateriaActualizada, err := q.UpdateMateria(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MateriaActualizada)
}

func deleteMateria(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	idStr := r.PathValue("id")

	parsed, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero de 32 bits válido", http.StatusBadRequest)
		return
	}

	idMateria := int32(parsed)

	err = q.DeleteMateria(ctx, idMateria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ==========================================
// 3. OPERACIONES DE TAREA (db.Tarea)
// ==========================================

func createTarea(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	var newTarea db.Tarea

	err := json.NewDecoder(r.Body).Decode(&newTarea)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newTarea.Titulo) == "" || newTarea.IDLista <= 0 {
		http.Error(w, "Todos los campos son obligatorios y no pueden contener solo espacios", http.StatusBadRequest)
		return
	}

	params := db.CreateTareaParams{
		IDLista:           newTarea.IDLista,
		Titulo:            newTarea.Titulo,
		FechaInicio:       newTarea.FechaInicio,
		FechaFinalizacion: newTarea.FechaFinalizacion,
		Descripcion:       newTarea.Descripcion,
	}

	tareaCreada, err := q.CreateTarea(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tareaCreada)
}

func getTareaById(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	idStr := r.PathValue("id")

	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de tarea inválido", http.StatusBadRequest)
		return
	}

	id := int32(id64)

	tarea, err := q.GetTareaByID(ctx, id)
	if err != nil {
		http.Error(w, "Tarea no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tarea)
}

func getListTareasByMateria(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	idStr := r.PathValue("id")

	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de materia inválido", http.StatusBadRequest)
		return
	}

	id := int32(id64)

	listaTareas, err := q.ListTareasByMateria(ctx, id)
	if err != nil {
		http.Error(w, "Materia para listar tareas no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listaTareas)
}

func getListTareasByRangoFechas(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	var tareaRango db.Tarea

	err := json.NewDecoder(r.Body).Decode(&tareaRango)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := db.ListTareasByRangoFechasParams{
		FechaFinalizacion:   tareaRango.FechaInicio,
		FechaFinalizacion_2: tareaRango.FechaFinalizacion,
	}

	listaTareas, err := q.ListTareasByRangoFechas(ctx, params)
	if err != nil {
		http.Error(w, "Tareas en rango de fechas no encontradas", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listaTareas)
}

func putUpdateTarea(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	var TareaRecibida db.Tarea

	err := json.NewDecoder(r.Body).Decode(&TareaRecibida)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(TareaRecibida.Titulo) == "" ||
		TareaRecibida.IDTarea <= 0 ||
		TareaRecibida.IDLista <= 0 {
		http.Error(w, "Todos los campos son obligatorios y no pueden contener solo espacios", http.StatusBadRequest)
		return
	}

	params := db.UpdateTareaParams{
		IDTarea:           TareaRecibida.IDTarea,
		Titulo:            TareaRecibida.Titulo,
		FechaInicio:       TareaRecibida.FechaInicio,
		FechaFinalizacion: TareaRecibida.FechaFinalizacion,
		Descripcion:       TareaRecibida.Descripcion,
		Completada:        TareaRecibida.Completada,
	}

	TareaActualizada, err := q.UpdateTarea(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TareaActualizada)
}

func deleteTarea(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	idStr := r.PathValue("id")

	parsed, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero de 32 bits válido", http.StatusBadRequest)
		return
	}

	idTarea := int32(parsed)

	err = q.DeleteTarea(ctx, idTarea)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ==========================================
// 4. OPERACIONES DE HORAS (db.HorasTrabajada)
// ==========================================

func createHorasTrabajadas(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	var newHoras db.HorasTrabajada

	err := json.NewDecoder(r.Body).Decode(&newHoras)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newHoras.IDTarea <= 0 {
		http.Error(w, "El ID de la tarea es obligatorio", http.StatusBadRequest)
		return
	}

	params := db.CreateHorasTrabajadasParams{
		IDTarea:         newHoras.IDTarea,
		TimestampInicio: newHoras.TimestampInicio,
		TimestampFin:    newHoras.TimestampFin,
	}

	horasCreadas, err := q.CreateHorasTrabajadas(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(horasCreadas)
}

func getListHorasByTarea(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	idStr := r.PathValue("id")

	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID de tarea inválido", http.StatusBadRequest)
		return
	}

	id := int32(id64)

	listaHoras, err := q.ListHorasByTarea(ctx, id)
	if err != nil {
		http.Error(w, "Tarea para listar horas no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listaHoras)
}

func getListHorasByRangoFechas(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	var horasRango db.HorasTrabajada

	err := json.NewDecoder(r.Body).Decode(&horasRango)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := db.ListHorasByRangoFechasParams{
		TimestampInicio:   horasRango.TimestampInicio,
		TimestampInicio_2: horasRango.TimestampFin,
	}

	listaHoras, err := q.ListHorasByRangoFechas(ctx, params)
	if err != nil {
		http.Error(w, "Horas en rango de fechas no encontradas", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listaHoras)
}

func putUpdateHorasTrabajadas(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	var HorasRecibidas db.HorasTrabajada

	err := json.NewDecoder(r.Body).Decode(&HorasRecibidas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if HorasRecibidas.IDHoras <= 0 || HorasRecibidas.IDTarea <= 0 {
		http.Error(w, "Todos los campos de ID son obligatorios", http.StatusBadRequest)
		return
	}

	params := db.UpdateHorasTrabajadasParams{
		IDHoras:         HorasRecibidas.IDHoras,
		TimestampInicio: HorasRecibidas.TimestampInicio,
		TimestampFin:    HorasRecibidas.TimestampFin,
	}

	HorasActualizadas, err := q.UpdateHorasTrabajadas(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HorasActualizadas)
}

func deleteHorasTrabajadas(w http.ResponseWriter, r *http.Request, q *db.Queries, ctx context.Context) {
	idStr := r.PathValue("id")

	parsed, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "El ID debe ser un número entero de 32 bits válido", http.StatusBadRequest)
		return
	}

	idHoras := int32(parsed)

	err = q.DeleteHorasTrabajadas(ctx, idHoras)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ==========================================
// FUNCIÓN MAIN
// ==========================================

func main() {
	// Conexion con la base de datos
	ctx := context.Background()

	connStr := "postgres://trabajo_web:carlos-luca-hermes@localhost:5432/go_tasks_db?sslmode=disable"

	conn, errOpen := abridDB(ctx, connStr)
	if errOpen != nil {
		fmt.Printf("error crítico al conectar a la base de datos: %v", errOpen)
		return
	}
	defer conn.Close(ctx)

	// Instanciar objeto de la base de datos (userRepository)
	queries := db.New(conn)

	// Implementacion del servidor de archivos
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	http.Handle("/", fileServer)

	// Handlers que conectan las operaciones crud
	// Usuarios[cite: 1]
	http.HandleFunc("POST /usuarios", func(w http.ResponseWriter, r *http.Request) {
		createUsuario(w, r, queries, ctx)
	})
	http.HandleFunc("GET /usuarios/email/{email}", func(w http.ResponseWriter, r *http.Request) {
		getUserByEmail(w, r, queries, ctx)
	})
	http.HandleFunc("GET /usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		getUserById(w, r, queries, ctx)
	})

	// Materias[cite: 1]
	http.HandleFunc("POST /materias", func(w http.ResponseWriter, r *http.Request) {
		createMateria(w, r, queries, ctx)
	})
	http.HandleFunc("GET /materias/{id}", func(w http.ResponseWriter, r *http.Request) {
		getListMateriasByUsuario(w, r, queries, ctx)
	})
	http.HandleFunc("PUT /materias/{id}", func(w http.ResponseWriter, r *http.Request) {
		putUpdateMateria(w, r, queries, ctx)
	})
	http.HandleFunc("DELETE /materias/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteMateria(w, r, queries, ctx)
	})

	// Tareas[cite: 1]
	http.HandleFunc("POST /tareas", func(w http.ResponseWriter, r *http.Request) {
		createTarea(w, r, queries, ctx)
	})
	http.HandleFunc("GET /tareas/{id}", func(w http.ResponseWriter, r *http.Request) {
		getTareaById(w, r, queries, ctx)
	})
	http.HandleFunc("GET /tareas/materia/{id}", func(w http.ResponseWriter, r *http.Request) {
		getListTareasByMateria(w, r, queries, ctx)
	})
	http.HandleFunc("GET /tareas/fechas", func(w http.ResponseWriter, r *http.Request) {
		getListTareasByRangoFechas(w, r, queries, ctx)
	})
	http.HandleFunc("PUT /tareas/{id}", func(w http.ResponseWriter, r *http.Request) {
		putUpdateTarea(w, r, queries, ctx)
	})
	http.HandleFunc("DELETE /tareas/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteTarea(w, r, queries, ctx)
	})

	// Horas[cite: 1]
	http.HandleFunc("POST /horas", func(w http.ResponseWriter, r *http.Request) {
		createHorasTrabajadas(w, r, queries, ctx)
	})
	http.HandleFunc("GET /horas/{id}", func(w http.ResponseWriter, r *http.Request) {
		getListHorasByTarea(w, r, queries, ctx)
	})
	http.HandleFunc("GET /horas/fechas", func(w http.ResponseWriter, r *http.Request) {
		getListHorasByRangoFechas(w, r, queries, ctx)
	})
	http.HandleFunc("PUT /horas/{id}", func(w http.ResponseWriter, r *http.Request) {
		putUpdateHorasTrabajadas(w, r, queries, ctx)
	})
	http.HandleFunc("DELETE /horas/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteHorasTrabajadas(w, r, queries, ctx)
	})

	port := ":8080"
	fmt.Printf("Servidor ESTÁTICO escuchando en http://localhost%s\n", port)
	fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
