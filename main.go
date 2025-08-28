package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Response struct {
	Status string `json:"status"`
    Message string `json:"message"`
	Data string `json:"data,omitempty"`
}

var db *sql.DB


func getInfoAbout(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    response := Response{Message: "Hello, World!", Status: "success"}
    json.NewEncoder(w).Encode(response)
}

func getInfoFromDb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	count := 0

    row := db.QueryRow("SELECT count(*) as count from cars_csv;")
	err := row.Scan(&count)

	// Если ошибка 
    if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{Message: "Неизвестная ошибка сервера", Status: "error"}
		json.NewEncoder(w).Encode(response)
		fmt.Println(err)
		return
    }

	w.WriteHeader(http.StatusOK)
    response := Response{Message: strconv.Itoa(count), Status: "success"}
    json.NewEncoder(w).Encode(response)
	
}

func main() {
	config, err := loadConfig()

	if err != nil {
		fmt.Print(err)
		return
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=%s", config.DB.User, config.DB.Password, config.DB.Dbname, config.DB.Port)
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Нет подключения к базе данных!")
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = db.Ping()
    if err != nil {
		fmt.Println("Нет подключения к базе данных!")
		fmt.Println(err)
		return
    }

    http.HandleFunc("/about", getInfoAbout)
	http.HandleFunc("/info", getInfoFromDb)

	fmt.Println("Server is listening...")
    http.ListenAndServe(":8080", nil)
}