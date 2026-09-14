package main

import (
	"ServidorEntregable/db"
	"context"
	"testing"
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
}
