package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
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

// Serve HTML Page
func homePageHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("./index.html")
    if err != nil {
        http.Error(w, "Could not load HTML page: "+err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, todos)
}

// Create a new todo from the form submission
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        title := r.FormValue("title")
        currentID++
        todo := Todo{ID: currentID, Title: title, Done: false}
        todos = append(todos, todo)
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Toggle Done status of a todo
func toggleDoneHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid id parameter", http.StatusBadRequest)
        return
    }
    for i := range todos {
        if todos[i].ID == id {
            todos[i].Done = !todos[i].Done
            break
        }
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Delete a todo
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid id parameter", http.StatusBadRequest)
        return
    }
    for i, todo := range todos {
        if todo.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            break
        }
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
    http.HandleFunc("/", homePageHandler)         // Home page with form and todo list
    http.HandleFunc("/todo/create", createTodoHandler)  // Create a new todo
    http.HandleFunc("/todo/toggle", toggleDoneHandler)  // Toggle done status
    http.HandleFunc("/todo/delete", deleteTodoHandler)  // Delete a todo

    fmt.Println("Server running at http://localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}
