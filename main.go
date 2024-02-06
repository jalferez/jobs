package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Connecting to database...")
	var err error
	db, err = sql.Open("pgx", "postgres://johnalferez:204851@localhost:5432/jobs?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	//API router for students
	router.HandleFunc("/students", getStudents).Methods("GET")
	router.HandleFunc("/students/{id}", getStudent).Methods("GET")
	router.HandleFunc("/students", addStudent).Methods("POST")
	router.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
