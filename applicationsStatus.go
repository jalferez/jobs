package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ApplicationStatus struct {
	Application_id              int       `json:"application_id"`
	Student_id                  int       `json:"student_id"`
	Position_id                 int       `json:"position_id"`
	Supervisor_id               int       `json:"supervisor_id"`
	Status                      string    `json:"status"`
	Application_submission_date time.Time `json:"application_submission_date"`
	Feedback                    string    `json:"feedback"`
}

func getApplicationsStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Application Status")
	rows, err := db.Query(`SELECT 
	application_id, 
	student_id, 
	position_id, 
	supervisor_id, 
	status, 
	application_submission_date, 
	feedback 
	FROM ApplicationStatus ORDER BY application_id ASC`)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	applicationStatus := []*ApplicationStatus{}
	for rows.Next() {
		application := &ApplicationStatus{}
		err := rows.Scan(&application.Application_id, &application.Student_id, &application.Position_id, &application.Supervisor_id, &application.Status, &application.Application_submission_date, &application.Feedback)
		if err != nil {
			log.Println("Error scanning rows:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		applicationStatus = append(applicationStatus, application)
	}

	result, err := json.Marshal(applicationStatus)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getApplicationStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Application")
	vars := mux.Vars(r)
	application_id := vars["id"]
	rows, err := db.Query(`SELECT 
	application_id, 
	student_id, 
	position_id, 
	supervisor_id, 
	status, 
	application_submission_date, 
	feedback 
	FROM ApplicationStatus WHERE application_id = $1`, application_id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	application := &ApplicationStatus{}
	for rows.Next() {
		err := rows.Scan(&application.Application_id, &application.Student_id, &application.Position_id, &application.Supervisor_id, &application.Status, &application.Application_submission_date, &application.Feedback)
		if err != nil {
			log.Println("Error scanning rows:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	result, err := json.Marshal(application)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func addApplicationStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding Application")
	application := &ApplicationStatus{}
	err := json.NewDecoder(r.Body).Decode(application)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`INSERT INTO ApplicationStatus 
	(student_id, position_id, supervisor_id, status, application_submission_date, feedback) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
		application.Student_id, application.Position_id, application.Supervisor_id, application.Status, application.Application_submission_date, application.Feedback)
	if err != nil {
		log.Println("Error inserting into database:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating Application")
	application := &ApplicationStatus{}
	err := json.NewDecoder(r.Body).Decode(application)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`UPDATE ApplicationStatus 
	SET student_id = $1, position_id = $2, supervisor_id = $3, status = $4, application_submission_date = $5, feedback = $6 
	WHERE application_id = $7`,
		application.Student_id, application.Position_id, application.Supervisor_id, application.Status, application.Application_submission_date, application.Feedback, application.Application_id)
	if err != nil {
		log.Println("Error updating row:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteApplicationStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting Application")
	vars := mux.Vars(r)
	application_id := vars["id"]
	_, err := db.Exec(`DELETE FROM ApplicationStatus WHERE application_id = $1`, application_id)
	if err != nil {
		log.Println("Error deleting row:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
