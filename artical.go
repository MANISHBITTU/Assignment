package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type Artical struct {
	ID          int64  `json:"Id"`
	TITLE       string `json:"title"`
	DESCRIPTION string `json:"description"`
	CONTENT     string `json:"content"`
}

func main() {

	db, err = gorm.Open("mysql", "manish:manish@tcp(localhost:3306)/student?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Connetion Failed")
	} else {
		log.Println("Connetion established")
	}

	db.AutoMigrate(&Artical{})
	handleRequest()

}

func handleRequest() {
	log.Println("Server Started:")
	log.Println("press ctrl and c to quit the server")
	myrouter := mux.NewRouter().StrictSlash(true)
	myrouter.HandleFunc("/new-artical", createNewArtical).Methods("GET")
	myrouter.HandleFunc("/all-articals", returnAllArticals).Methods("GET")
	myrouter.HandleFunc("/artical/{id}", returnSingleArtical).Methods("GET")
	myrouter.HandleFunc("/update-artical/{id}/{title}", updateArtical).Methods("GET")
	myrouter.HandleFunc("/delete-artical/{id}", deleteArtical).Methods("GET")
	log.Fatal(http.ListenAndServe(":8085", myrouter))
}

func createNewArtical(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		var artical Artical
		er := json.Unmarshal(reqBody, &artical)
		if er == nil {
			db.Create(&artical)
			fmt.Println("Endpoint Hit: Creating New Booking")
			json.NewEncoder(w).Encode(artical)
		} else {
			fmt.Println(er)
		}
	}
}

func returnTotalArticals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	articals := []Artical{}
	db.Raw("Select * from aricals").Scan(&articals)
	json.NewEncoder(w).Encode(articals)
}

func returnSingleBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	artical := []Artical{}
	s, err := strconv.Atoi(key)
	if err == nil {
		db.Raw("Select * from articals where id= ?", s).Scan(&artical)
		json.NewEncoder(w).Encode(artical)
	}

}

func updateBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key1 := vars["id"]
	key2 := vars["user"]
	articals := []Artical{}
	s, err := strconv.Atoi(key1)
	if err == nil {
		db.Model(&articals).Where("id = ?", s).Update("title", key2)
		fmt.Println("Updated")
	} else {
		fmt.Print(err)
	}

}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	articals := []Artical{}
	s, err := strconv.Atoi(key)
	if err == nil {
		db.Where("id = ?", s).Delete(&articals)
		fmt.Println("Deleted")
	} else {
		fmt.Print(err)
	}
}
