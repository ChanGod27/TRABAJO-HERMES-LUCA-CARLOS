package main

import (
	"ServidorEntregable/db"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestQueries_CRUD(t *testing.T) {
	//Al usar el drive pgx, hay que definir un context
	//Para las interacciones con la api
	ctx := context.Background()

	//Definimos el string con los parametros de conexion de la base de datos
	connStr := "postgres://trabajo_web:carlos-luca-hermes@localhost:5432/go_tasks_db?sslmode=disable"

	//Abrimos la base de datos
	conn, errOpen := abridDB(ctx, connStr)
	if errOpen != nil {
		t.Fatal(errOpen)
	}
	//Cuando termine de ejecutar el metodo, se cierra la conexion con la base de datos
	defer conn.Close(ctx)

	//Instanciamos el objeto de queries
	queries := db.New(conn)

	//TESTS DE OPERACIONES CRUD
	t.Log("Iniciando pruebas de operaciones CRUD")

	// =========================================================================
	// 1. TESTS CRUD: USUARIO
	// =========================================================================

	// 1.1 Create User
	userParams := db.CreateUserParams{
		Nombre:     "Hermes",
		Apellido:   "Truelsegaard",
		Email:      "hermes@example.com",
		Contrasena: "hash_super_secreto_123",
	}

	userCreado, err := queries.CreateUser(ctx, userParams)
	if err != nil {
		t.Fatalf("CreateUser falló: %v", err)
	}

	if userCreado.IDUsuario == 0 {
		t.Fatalf("CreateUser retornó un ID inválido (0)")
	}

	t.Logf("Usuario creado exitosamente con ID: %d", userCreado.IDUsuario)

	// 1.2 Get User By ID
	userPorID, err := queries.GetUserByID(ctx, userCreado.IDUsuario)
	if err != nil {
		t.Fatalf("GetUserByID falló: %v", err)
	}
	if userPorID.Email != userParams.Email {
		t.Fatalf("GetUserByID inconsistente: esperado %s, obtenido %s", userParams.Email, userPorID.Email)
	}

	// 1.3 Get User By Email
	userPorEmail, err := queries.GetUserByEmail(ctx, userParams.Email)
	if err != nil {
		t.Fatalf("GetUserByEmail falló: %v", err)
	}
	if userPorEmail.IDUsuario != userCreado.IDUsuario {
		t.Fatalf("GetUserByEmail inconsistente: esperado ID %d, obtenido %d", userCreado.IDUsuario, userPorEmail.IDUsuario)
	}

	// =========================================================================
	// 2. TESTS CRUD: MATERIA (LISTA_TAREA)
	// =========================================================================
	t.Log("Iniciando pruebas de LISTA_TAREA (MATERIA)...")

	// 2.1 Create Materia
	materiaParams := db.CreateMateriaParams{
		IDUsuario: userCreado.IDUsuario,
		Titulo:    "Programación Web",
	}
	materiaCreada, err := queries.CreateMateria(ctx, materiaParams)
	if err != nil {
		t.Fatalf("CreateMateria falló: %v", err)
	}
	if materiaCreada.IDLista == 0 {
		t.Fatalf("CreateMateria retornó un ID inválido (0)")
	}
	t.Logf("Materia creada con ID: %d", materiaCreada.IDLista)

	// 2.2 List Materias By Usuario
	materias, err := queries.ListMateriasByUsuario(ctx, userCreado.IDUsuario)
	if err != nil {
		t.Fatalf("ListMateriasByUsuario falló: %v", err)
	}
	if len(materias) != 1 {
		t.Fatalf("ListMateriasByUsuario esperaba 1 elemento, obtuvo %d", len(materias))
	}

	// 2.3 Update Materia
	updateMateriaParams := db.UpdateMateriaParams{
		IDLista: materiaCreada.IDLista,
		Titulo:  "Programación Web II",
	}
	materiaActualizada, err := queries.UpdateMateria(ctx, updateMateriaParams)
	if err != nil {
		t.Fatalf("UpdateMateria falló: %v", err)
	}
	if materiaActualizada.Titulo != updateMateriaParams.Titulo {
		t.Fatalf("UpdateMateria no actualizó el título: esperado %s, obtenido %s", updateMateriaParams.Titulo, materiaActualizada.Titulo)
	}

	// =========================================================================
	// 3. TESTS CRUD: TAREA
	// =========================================================================
	t.Log("Iniciando pruebas de TAREA...")

	ahora := time.Now().Truncate(time.Microsecond)
	vencimiento := ahora.Add(48 * time.Hour)

	// 3.1 Create Tarea
	tareaParams := db.CreateTareaParams{
		IDLista:           materiaCreada.IDLista,
		Titulo:            "Entrega de TP2",
		FechaInicio:       pgtype.Timestamp{Time: ahora, Valid: true},
		FechaFinalizacion: pgtype.Timestamp{Time: vencimiento, Valid: true},
		Descripcion:       pgtype.Text{String: "Implementar sqlc y testing", Valid: true},
	}
	tareaCreada, err := queries.CreateTarea(ctx, tareaParams)
	if err != nil {
		t.Fatalf("CreateTarea falló: %v", err)
	}
	if tareaCreada.IDTarea == 0 {
		t.Fatalf("CreateTarea retornó un ID inválido (0)")
	}
	t.Logf("Tarea creada con ID: %d", tareaCreada.IDTarea)

	// 3.2 Get Tarea By ID
	tareaPorID, err := queries.GetTareaByID(ctx, tareaCreada.IDTarea)
	if err != nil {
		t.Fatalf("GetTareaByID falló: %v", err)
	}
	if tareaPorID.Titulo != tareaParams.Titulo {
		t.Fatalf("GetTareaByID inconsistente: esperado %s, obtenido %s", tareaParams.Titulo, tareaPorID.Titulo)
	}

	// 3.3 List Tareas By Materia
	tareas, err := queries.ListTareasByMateria(ctx, materiaCreada.IDLista)
	if err != nil {
		t.Fatalf("ListTareasByMateria falló: %v", err)
	}
	if len(tareas) != 1 {
		t.Fatalf("ListTareasByMateria esperaba 1 tarea, obtuvo %d", len(tareas))
	}

	// 3.4 List Tareas By Rango de Fechas
	tareasRango, err := queries.ListTareasByRangoFechas(ctx, db.ListTareasByRangoFechasParams{
		FechaFinalizacion:   pgtype.Timestamp{Time: ahora.Add(-1 * time.Hour), Valid: true},
		FechaFinalizacion_2: pgtype.Timestamp{Time: vencimiento.Add(1 * time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("ListTareasByRangoFechas falló: %v", err)
	}
	if len(tareasRango) == 0 {
		t.Fatalf("ListTareasByRangoFechas esperaba al menos 1 resultado, obtuvo 0")
	}

	// 3.5 Update Tarea
	updateTareaParams := db.UpdateTareaParams{
		IDTarea:           tareaCreada.IDTarea,
		Titulo:            "Entrega de TP2 - Finalizada",
		FechaInicio:       tareaCreada.FechaInicio,
		FechaFinalizacion: tareaCreada.FechaFinalizacion,
		Descripcion:       pgtype.Text{String: "Entrega completada con tests", Valid: true},
		Completada:        true,
	}
	tareaActualizada, err := queries.UpdateTarea(ctx, updateTareaParams)
	if err != nil {
		t.Fatalf("UpdateTarea falló: %v", err)
	}
	if !tareaActualizada.Completada || tareaActualizada.Titulo != updateTareaParams.Titulo {
		t.Fatalf("UpdateTarea no persistió los cambios adecuadamente")
	}

	// =========================================================================
	// 4. TESTS CRUD: HORAS TRABAJADAS (CRONÓMETRO)
	// =========================================================================
	t.Log("Iniciando pruebas de HORAS TRABAJADAS...")

	inicioSesion := ahora.Add(1 * time.Hour)
	finSesion := ahora.Add(3 * time.Hour)

	// 4.1 Create Horas Trabajadas
	horasParams := db.CreateHorasTrabajadasParams{
		IDTarea:         tareaCreada.IDTarea,
		TimestampInicio: pgtype.Timestamp{Time: inicioSesion, Valid: true},
		TimestampFin:    pgtype.Timestamp{Time: finSesion, Valid: true},
	}
	horasCreadas, err := queries.CreateHorasTrabajadas(ctx, horasParams)
	if err != nil {
		t.Fatalf("CreateHorasTrabajadas falló: %v", err)
	}
	if horasCreadas.IDHoras == 0 {
		t.Fatalf("CreateHorasTrabajadas retornó un ID inválido (0)")
	}
	t.Logf("Horas trabajadas registradas con ID: %d", horasCreadas.IDHoras)

	// 4.2 List Horas By Tarea
	horasPorTarea, err := queries.ListHorasByTarea(ctx, tareaCreada.IDTarea)
	if err != nil {
		t.Fatalf("ListHorasByTarea falló: %v", err)
	}
	if len(horasPorTarea) != 1 {
		t.Fatalf("ListHorasByTarea esperaba 1 registro, obtuvo %d", len(horasPorTarea))
	}

	// 4.3 List Horas By Rango de Fechas
	horasRango, err := queries.ListHorasByRangoFechas(ctx, db.ListHorasByRangoFechasParams{
		TimestampInicio:   pgtype.Timestamp{Time: ahora, Valid: true},
		TimestampInicio_2: pgtype.Timestamp{Time: finSesion.Add(1 * time.Hour), Valid: true},
	})
	if err != nil {
		t.Fatalf("ListHorasByRangoFechas falló: %v", err)
	}
	if len(horasRango) == 0 {
		t.Fatalf("ListHorasByRangoFechas esperaba al menos 1 resultado, obtuvo 0")
	}

	// 4.4 Update Horas Trabajadas
	finSesionCorregido := finSesion.Add(30 * time.Minute).Truncate(time.Second)
	updateHorasParams := db.UpdateHorasTrabajadasParams{
		IDHoras:         horasCreadas.IDHoras,
		TimestampInicio: horasCreadas.TimestampInicio,
		TimestampFin:    pgtype.Timestamp{Time: finSesionCorregido, Valid: true},
	}
	horasActualizadas, err := queries.UpdateHorasTrabajadas(ctx, updateHorasParams)
	if err != nil {
		t.Fatalf("UpdateHorasTrabajadas falló: %v", err)
	}
	if horasActualizadas.TimestampFin.Time.Unix() != finSesionCorregido.Unix() {
		t.Fatalf("UpdateHorasTrabajadas inconsistente: esperado %v, obtenido %v",
			finSesionCorregido, horasActualizadas.TimestampFin.Time)
	}

	// =========================================================================
	// 5. BORRADO EN CASCADA / LIMPIEZA INVERSA (DELETE)
	// =========================================================================
	t.Log("Iniciando eliminación secuencial de registros...")

	// 5.1 Delete Horas Trabajadas
	if err := queries.DeleteHorasTrabajadas(ctx, horasCreadas.IDHoras); err != nil {
		t.Fatalf("DeleteHorasTrabajadas falló: %v", err)
	}
	horasCheck, err := queries.ListHorasByTarea(ctx, tareaCreada.IDTarea)
	if err != nil {
		t.Fatalf("Error consultando horas tras borrado: %v", err)
	}
	if len(horasCheck) != 0 {
		t.Fatalf("DeleteHorasTrabajadas no eliminó el registro: quedan %d items", len(horasCheck))
	}

	// 5.2 Delete Tarea
	if err := queries.DeleteTarea(ctx, tareaCreada.IDTarea); err != nil {
		t.Fatalf("DeleteTarea falló: %v", err)
	}
	_, err = queries.GetTareaByID(ctx, tareaCreada.IDTarea)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("GetTareaByID esperaba pgx.ErrNoRows tras el borrado, obtuvo: %v", err)
	}

	// 5.3 Delete Materia
	if err := queries.DeleteMateria(ctx, materiaCreada.IDLista); err != nil {
		t.Fatalf("DeleteMateria falló: %v", err)
	}
	materiasCheck, err := queries.ListMateriasByUsuario(ctx, userCreado.IDUsuario)
	if err != nil {
		t.Fatalf("Error consultando materias tras borrado: %v", err)
	}
	if len(materiasCheck) != 0 {
		t.Fatalf("DeleteMateria no eliminó el registro: quedan %d items", len(materiasCheck))
	}

	t.Log("Todos los tests de operaciones CRUD pasaron exitosamente.")

}
