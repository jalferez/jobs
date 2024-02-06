package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Supervisor struct {
	Supervisor_id   int    `json:"supervisor_id"`
	First_name      string `json:"first_name"`
	Last_name       string `json:"last_name"`
	Email           string `json:"email"`
	Phone_number    string `json:"phone_number"`
	Department_name string `json:"department_name"`
}

func getSupervisors(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Supervisors")
	rows, err := db.Query(`SELECT 
	supervisor_id, 
	first_name, 
	last_name, 
	email, 
	phone_number, 
	department_name
	FROM supervisors ORDER BY supervisor_id ASC`)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	supervisors := []*Supervisor{}
	for rows.Next() {
		supervisor := &Supervisor{}
		err := rows.Scan(&supervisor.Supervisor_id, &supervisor.First_name, &supervisor.Last_name, &supervisor.Email, &supervisor.Phone_number, &supervisor.Department_name)
		if err != nil {
			log.Println("Error scanning rows:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		supervisors = append(supervisors, supervisor)
	}

	result, err := json.Marshal(supervisors)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getSupervisor(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Supervisor")
	params := mux.Vars(r)
	id := params["id"]
	row := db.QueryRow(`SELECT 
	supervisor_id, 
	first_name, 
	last_name, 
	email, 
	phone_number, 
	department_name
	FROM supervisors WHERE supervisor_id = $1`, id)

	supervisor := &Supervisor{}
	err := row.Scan(&supervisor.Supervisor_id, &supervisor.First_name, &supervisor.Last_name, &supervisor.Email, &supervisor.Phone_number, &supervisor.Department_name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(supervisor)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func addSupervisor(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding Supervisor")
	var supervisor Supervisor
	err := json.NewDecoder(r.Body).Decode(&supervisor)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`INSERT INTO supervisors 
	(first_name, last_name, email, phone_number, department_name) 
	VALUES ($1, $2, $3, $4, $5)`,
		supervisor.First_name, supervisor.Last_name, supervisor.Email, supervisor.Phone_number, supervisor.Department_name)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateSupervisor(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating Supervisor")
	var supervisor Supervisor
	err := json.NewDecoder(r.Body).Decode(&supervisor)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`UPDATE supervisors 
	SET first_name = $1, last_name = $2, email = $3, phone_number = $4, department_name = $5
	WHERE supervisor_id = $6`,
		supervisor.First_name, supervisor.Last_name, supervisor.Email, supervisor.Phone_number, supervisor.Department_name, supervisor.Supervisor_id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteSupervisor(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting Supervisor")
	params := mux.Vars(r)
	id := params["id"]
	_, err := db.Exec(`DELETE FROM supervisors WHERE supervisor_id = $1`, id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
