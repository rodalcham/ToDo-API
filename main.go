package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// Todo struct defines the fields for a to-do item.
type Todo struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

// Slice to store all Todo items
var todos []Todo
var currentID = 0

// Create adds a new todo to the list
func create(title string) Todo {
    currentID++
    todo := Todo{ID: currentID, Title: title, Done: false}
    todos = append(todos, todo)
    return todo
}

// Read returns all todos as JSON
func read(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

// Update marks a todo as done or updates its title
func update(id int, title string, done bool) bool {
    for i, todo := range todos {
        if todo.ID == id {
            todos[i].Title = title
            todos[i].Done = done
            return true
        }
    }
    return false
}

// Delete removes a todo from the list
func delete(id int) bool {
    for i, todo := range todos {
        if todo.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            return true
        }
    }
    return false
}

func main() {
    // Initialize the server and routes
    http.HandleFunc("/todos", read) // Route to get all todos as JSON

    // Adding some initial todos
    create("Learn Go")
    create("Implement REST API")
    
    fmt.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
