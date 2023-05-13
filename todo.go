package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

type Todolist struct {
	ID       string `json: "id"`
	Name     string `json: "name"`
	Duration int    `json: "duration"`
}

var todo []Todolist

// get all todo:
func getTodoLists(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// get a todo by id
func getTodoList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)

	for _, value := range todo {

		if value.ID == param["id"] {
			json.NewEncoder(w).Encode(value)
			break
		}
	}
}

// delete a todo:

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)

	for index, value := range todo {
		if value.ID == param["id"] {
			todo = append(todo[:index], todo[index+1:]...)
		}
	}
	json.NewEncoder(w).Encode(todo)
}

func deleteTodolist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todo = todo[:0]
	json.NewEncoder(w).Encode(todo)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTodo Todolist
	_ = json.NewDecoder(r.Body).Decode(&newTodo)
	newTodo.ID = strconv.Itoa(rand.Intn(100000))

	todo = append(todo, newTodo)

	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var updateTodo Todolist
	param := mux.Vars(r)

	for index, value := range todo {

		if value.ID == param["id"] {

			todo = append(todo[:index], todo[index+1:]...)
			_ = json.NewDecoder(r.Body).Decode(&updateTodo)
			updateTodo.ID = strconv.Itoa(rand.Intn(100000))

			todo = append(todo, updateTodo)
			json.NewEncoder(w).Encode(todo)

		}
	}
}

func sortTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sort.Slice(todo, func(i, j int) bool {
		return todo[i].Duration < todo[j].Duration
	})
}

func main() {

	todo = append(todo, Todolist{ID: "1", Name: "Cook Lunch", Duration: 45})
	todo = append(todo, Todolist{ID: "2", Name: "Eat Lunch", Duration: 15})
	r := mux.NewRouter()

	r.HandleFunc("/todo", getTodoLists).Methods("GET")
	r.HandleFunc("/todo/{id}", getTodoList).Methods("GET")
	r.HandleFunc("/todo/{duration}", sortTodo).Methods("GET")
	r.HandleFunc("/todo", createTodo).Methods("POST")
	r.HandleFunc("/todo/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")
	r.HandleFunc("/todo", deleteTodolist).Methods("DELETE")

	fmt.Printf("Starting the server at the local host 8000")

	log.Fatal(http.ListenAndServe(":8000", r))
}
