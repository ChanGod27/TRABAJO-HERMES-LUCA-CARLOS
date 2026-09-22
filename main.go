package main

import (
	"ServidorEntregable/db"
	"context"
	"fmt"
	"net/http"
)

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

}
