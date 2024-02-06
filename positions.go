package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Position struct {
	Position_id     int    `json:"position_id"`
	Department_name string `json:"department_name"`
	Position_name   string `json:"position_name"`
	Age_requirement int    `json:"age_requirement"`
	Semester        string `json:"semester"`
	Email           string `json:"email"`
	Phone_number    string `json:"phone_number"`
	Pay_rate        float64    `json:"pay_rate"`
}

func getPositions(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Postitions")
	rows, err := db.Query(`
	SELECT
	position_id,
	department_name,
	position_name,
	age_requirement,
	semester,
	email,
	phone_number,
	pay_rate
	FROM positions ORDER BY semester, position_id ASC`)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	positions := []*Position{}
	for rows.Next() {
		position := &Position{}
		err := rows.Scan(&position.Position_id, &position.Department_name, &position.Position_name, &position.Age_requirement, &position.Semester, &position.Email, &position.Phone_number, &position.Pay_rate)
		if err != nil {
			log.Println("Error scanning rows:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		positions = append(positions, position)
	}

	result, err := json.Marshal(positions)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func getPosition(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Position")
	params := mux.Vars(r)
	position_id := params["position_id"]
	row := db.QueryRow(`
	SELECT
	position_id,
	department_name,
	position_name,
	age_requirement,
	semester,
	email,
	phone_number,
	pay_rate
	FROM positions WHERE position_id = $1`, position_id)

	position := &Position{}
	err := row.Scan(&position.Position_id, &position.Department_name, &position.Position_name, &position.Age_requirement, &position.Semester, &position.Email, &position.Phone_number, &position.Pay_rate)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(position)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func addPosition(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding Position")
	position := &Position{}
	err := json.NewDecoder(r.Body).Decode(position)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
	INSERT INTO positions
	(department_name, position_name, age_requirement, semester, email, phone_number, pay_rate)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		position.Department_name, position.Position_name, position.Age_requirement, position.Semester, position.Email, position.Phone_number, position.Pay_rate)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updatePosition(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating Position")
	position := &Position{}
	err := json.NewDecoder(r.Body).Decode(position)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`
	UPDATE positions
	SET department_name = $1, position_name = $2, age_requirement = $3, semester = $4, email = $5, phone_number = $6, pay_rate = $7
	WHERE position_id = $8`,
		position.Department_name, position.Position_name, position.Age_requirement, position.Semester, position.Email, position.Phone_number, position.Pay_rate, position.Position_id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deletePosition(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting Position")
	params := mux.Vars(r)
	position_id := params["position_id"]

	_, err := db.Exec(`DELETE FROM positions WHERE position_id = $1`, position_id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
