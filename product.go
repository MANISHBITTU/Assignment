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

type Product struct {
	Product_id    int64   `json:"Product_id"`
	Store_id      int64   `json:"Store_id"`
	Store_name    string  `json:"name"`
	Phone         string  `json:"phone"`
	City          string  `json:"city"`
	Product_Price float64 `json:"Product_Price"`
	Product_name  string  `json:"product_name"`
}

func main() {

	db, err = gorm.Open("mysql", "manish:manish@tcp(localhost:3306)/student?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Connetion Failed")
	} else {
		log.Println("Connetion established")
	}

	db.AutoMigrate(&Product{})
	handleRequest()

}

func handleRequest() {
	log.Println("Server Started:")
	log.Println("press ctrl and c to quit the server")
	myrouter := mux.NewRouter().StrictSlash(true)
	myrouter.HandleFunc("/new-product", createNewProduct).Methods("GET")
	myrouter.HandleFunc("/all-products", returnAllProducts).Methods("GET")
	myrouter.HandleFunc("/product/{id}", returnSingleProduct).Methods("GET")
	myrouter.HandleFunc("/update-product/{pid}/{sid}", updateProduct).Methods("GET")
	myrouter.HandleFunc("/delete-product/{id}", deleteProduct).Methods("GET")
	log.Fatal(http.ListenAndServe(":8085", myrouter))
}

func createNewProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		var product Product
		er := json.Unmarshal(reqBody, &product)
		if er == nil {
			db.Create(&product)
			fmt.Println("Endpoint Hit: Creating New Booking")
			json.NewEncoder(w).Encode(product)
		} else {
			fmt.Println(er)
		}
	}
}

func returnAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/json")
	products := []Product{}
	db.Raw("Select * from products").Scan(&products)
	json.NewEncoder(w).Encode(products)
}

func returnSingleProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	product := []Product{}
	s, err := strconv.Atoi(key)
	if err == nil {
		db.Raw("Select * from products where id= ?", s).Scan(&product)
		json.NewEncoder(w).Encode(product)
	}

}

func updateBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key1 := vars["pid"]
	key2 := vars["sid"]
	products := []Product{}
	s1, err := strconv.Atoi(key1)
	s2, err := strconv.Atoi(key2)
	if err == nil {
		db.Model(&articals).Where("product_id = ?", s1).Update("store_id", s2)
		fmt.Println("Updated")
	} else {
		fmt.Print(err)
	}

}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	products := []Product{}
	s, err := strconv.Atoi(key)
	if err == nil {
		db.Where("id = ?", s).Delete(&products)
		fmt.Println("Deleted")
	} else {
		fmt.Print(err)
	}
}
