package main

import (
	"encoding/json"

	"log"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"

	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Driver struct {
	gorm.Model

	Name string

	License string

	Cars []Car
}

type Car struct {
	gorm.Model

	Year int

	Make string

	ModelName string

	DriverID int
}

var db *gorm.DB

var err error

func main() {

	router := mux.NewRouter()

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=text sslmode=disable password=sandy5353")

	if err != nil {

		panic("failed to connect database")

	}

	defer db.Close()

	db.AutoMigrate(&Driver{})

	db.AutoMigrate(&Car{})

	router.HandleFunc("/cars", GetCars).Methods("GET")

	router.HandleFunc("/cars/{id}", GetCar).Methods("GET")

	router.HandleFunc("/drivers/{id}", GetDriver).Methods("GET")

	router.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")

	router.HandleFunc("/addcar", AddCar).Methods("POST")

	router.HandleFunc("/adddriver", AddDriver).Methods("POST")

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))

}

func GetCars(w http.ResponseWriter, r *http.Request) {

	var cars []Car

	db.Find(&cars)

	json.NewEncoder(w).Encode(&cars)

}

func GetCar(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var car Car

	db.First(&car, params["id"])

	json.NewEncoder(w).Encode(&car)

}

func GetDriver(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var driver Driver

	var cars []Car

	db.First(&driver, params["id"])

	db.Model(&driver).Related(&cars)

	driver.Cars = cars

	json.NewEncoder(w).Encode(&driver)

}

func DeleteCar(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var car Car

	db.First(&car, params["id"])

	db.Delete(&car)

	var cars []Car

	db.Find(&cars)

	json.NewEncoder(w).Encode(&cars)

}

func AddCar(w http.ResponseWriter, r *http.Request) {
	var car Car

	json.NewDecoder(r.Body).Decode(&car)

	db.Create(&car)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(car)
}

func AddDriver(w http.ResponseWriter, r *http.Request) {
	var driver Driver

	json.NewDecoder(r.Body).Decode(&driver)

	db.Create(&driver)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(driver)
}
