package main

import (
	"database/sql"
	"log"
	"net/http"
	"github.com/gorilla/handlers"
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

	log.Println("Creating new router...")
	router := mux.NewRouter()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	router.Use(corsHandler)


	log.Println("Registering routes...")
	//API routes for students
	router.HandleFunc("/api/v1/students", getStudents).Methods("GET")
	router.HandleFunc("/api/v1/students/{id}", getStudent).Methods("GET")
	router.HandleFunc("/api/v1/students", addStudent).Methods("POST")
	router.HandleFunc("/api/v1/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/api/v1/students/{id}", deleteStudent).Methods("DELETE")

	//API routes for positions
	router.HandleFunc("/api/v1/positions", getPositions).Methods("GET")
	router.HandleFunc("/api/v1/positions/{id}", getPosition).Methods("GET")
	router.HandleFunc("/api/v1/positions", addPosition).Methods("POST")
	router.HandleFunc("/api/v1/positions/{id}", updatePosition).Methods("PUT")
	router.HandleFunc("/api/v1/positions/{id}", deletePosition).Methods("DELETE")

	//API routes for supervisors
	router.HandleFunc("/api/v1/supervisors", getSupervisors).Methods("GET")
	router.HandleFunc("/api/v1/supervisors/{id}", getSupervisor).Methods("GET")
	router.HandleFunc("/api/v1/supervisors", addSupervisor).Methods("POST")
	router.HandleFunc("/api/v1/supervisors/{id}", updateSupervisor).Methods("PUT")
	router.HandleFunc("/api/v1/supervisors/{id}", deleteSupervisor).Methods("DELETE")

	//API routes for application status
	router.HandleFunc("/api/v1/applicationStatus", getApplicationsStatus).Methods("GET")
	router.HandleFunc("/api/v1/applicationStatus/{id}", getApplicationStatus).Methods("GET")
	router.HandleFunc("/api/v1/applicationStatus", addApplicationStatus).Methods("POST")
	router.HandleFunc("/api/v1/applicationStatus/{id}", updateApplicationStatus).Methods("PUT")
	router.HandleFunc("/api/v1/applicationStatus/{id}", deleteApplicationStatus).Methods("DELETE")

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
