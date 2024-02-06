package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Student struct {
	Student_id    int    `json:"student_id"`
	First_name    string `json:"first_name"`
	Last_name     string `json:"last_name"`
	Date_of_birth string `json:"date_of_birth"`
	Age           int    `json:"age"`
	Phone_number  string `json:"phone_number"`
	Email         string `json:"email"`
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Students")
	rows, err := db.Query(`SELECT 
	student_id, 
	first_name, 
	last_name, 
	date_of_birth, 
	age, 
	phone_number, 
	email 
	FROM students ORDER BY student_id ASC`)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	students := []*Student{}
	for rows.Next() {
		student := &Student{}
		err := rows.Scan(&student.Student_id, &student.First_name, &student.Last_name, &student.Date_of_birth, &student.Age, &student.Phone_number, &student.Email)
		if err != nil {
			log.Println("Error scanning rows:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	result, err := json.Marshal(students)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Students")
	vars := mux.Vars(r)
	id := vars["student_id"]
	if id == "" {
		log.Println("Missing student ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	row := db.QueryRow(`SELECT 
		student_id,
		first_name,
		last_name,
		date_of_birth,
		age,
		phone_number,
		email
		FROM students WHERE student_id = $1`,
		id)

	student := &Student{}
	err := row.Scan(&student.Student_id, &student.First_name, &student.Last_name, &student.Date_of_birth, &student.Age, &student.Phone_number, &student.Email)
	if err != nil {
		log.Println("Error scanning row:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(student)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func addStudent(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding Student")
	student := &Student{}
	err := json.NewDecoder(r.Body).Decode(student)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`INSERT INTO students 
		(first_name, last_name, date_of_birth, age, phone_number, email) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		student.First_name, student.Last_name, student.Date_of_birth, student.Age, student.Phone_number, student.Email)
	if err != nil {
		log.Println("Error inserting row:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating Student")
	student := &Student{}
	err := json.NewDecoder(r.Body).Decode(student)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`UPDATE students 
		SET first_name = $1, last_name = $2, date_of_birth = $3, age = $4, phone_number = $5, email = $6 
		WHERE student_id = $7`,
		student.First_name, student.Last_name, student.Date_of_birth, student.Age, student.Phone_number, student.Email, student.Student_id)
	if err != nil {
		log.Println("Error updating row:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting Student")
	vars := mux.Vars(r)
	id := vars["student_id"]
	if id == "" {
		log.Println("Missing student ID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`DELETE FROM students WHERE student_id = $1`, id)
	if err != nil {
		log.Println("Error deleting row:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
