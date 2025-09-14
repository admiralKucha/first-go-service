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
	Data interface{} `json:"data,omitempty"`
}

type CarSmallInfo struct {
	Id string `json:"id"`
    Name string `json:"message"`
}

type Counts struct {
	Name string `json:"string"`
    Count int `json:"int"`
}


var db *sql.DB


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

func getAllCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := make([]CarSmallInfo, 0)
	var offset int 

    // Смотрим какое смещение
    offsetStr := r.URL.Query().Get("offset")
    if offsetStr == "" {
        offset = 0 
    } else {
        var err error
        offset, err = strconv.Atoi(offsetStr)
        if err != nil {
            offset = 0
        } else {
			offset = offset - (offset % 10)
		}
    }

	rows, err := db.Query(
		"SELECT id, name FROM unique_summary_cars " +
		"LIMIT 10 " +
		"OFFSET $1;",
		offset,
	)

	// Если ошибка 
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{Message: "Неизвестная ошибка сервера", Status: "error"}
		json.NewEncoder(w).Encode(response)
		fmt.Println(err)
		return
    }

	// Обрабатываем полученный результат
	defer rows.Close()
    for rows.Next(){
        p := CarSmallInfo{}
        err := rows.Scan(&p.Id, &p.Name)

		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			response := Response{Message: "Неизвестная ошибка сервера", Status: "error"}
			json.NewEncoder(w).Encode(response)
			fmt.Println(err)
			return
		}

        data = append(data, p)
    }

	w.WriteHeader(http.StatusOK)
    response := Response{Message: "Получен список автомобилей", Status: "success", Data: data}
    json.NewEncoder(w).Encode(response)

}

func getAllYearsCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := make([]Counts, 0)

	rows, err := db.Query("SELECT date_of_publication_year, count(*) FROM unique_summary_cars GROUP BY date_of_publication_year;")

	// Если ошибка 
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{Message: "Неизвестная ошибка сервера", Status: "error"}
		json.NewEncoder(w).Encode(response)
		fmt.Println(err)
		return
    }

	// Обрабатываем полученный результат
	defer rows.Close()
    for rows.Next(){
        p := Counts{}
        err := rows.Scan(&p.Name, &p.Count)

		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			response := Response{Message: "Неизвестная ошибка сервера", Status: "error"}
			json.NewEncoder(w).Encode(response)
			return
		}

        data = append(data, p)
    }

	w.WriteHeader(http.StatusOK)
    response := Response{Message: "Получена информация о количестве машин годов", Status: "success", Data: data}
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

	http.HandleFunc("/info", getInfoFromDb)
	http.HandleFunc("/", getAllCars)
	http.HandleFunc("/years", getAllYearsCount)

	fmt.Println("Server is listening...")
    http.ListenAndServe(":8080", nil)
}