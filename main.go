package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
    "time"
)

// Todo struct defines the fields for a to-do item.
type Todo struct {
    ID    int
    Title string
    Done  bool
    Day   time.Time // Each Todo is associated with a specific day
}

var todos []Todo    // Slice to store all todos
var currentID = 0   // ID counter

// Struct to hold a day's date and its todos for the week view
type Day struct {
    Date  time.Time
    Todos []Todo
}

// Get the current week's days starting from today
func getCurrentWeek() []Day {
    var week []Day
    today := time.Now()
    startOfWeek := today.AddDate(0, 0, -int(today.Weekday()-time.Monday))

    for i := 0; i < 7; i++ {
        day := startOfWeek.AddDate(0, 0, i)
        week = append(week, Day{
            Date:  day,
            Todos: getTodosForDay(day),
        })
    }
    return week
}

// Retrieve todos for a specific day
func getTodosForDay(day time.Time) []Todo {
    var dayTodos []Todo
    for _, todo := range todos {
        if todo.Day.Year() == day.Year() && todo.Day.YearDay() == day.YearDay() {
            dayTodos = append(dayTodos, todo)
        }
    }
    return dayTodos
}

// Serve HTML Page with the weekly view
func homePageHandler(w http.ResponseWriter, r *http.Request) {
    week := getCurrentWeek()
    tmpl, err := template.ParseFiles("./index.html")
    if err != nil {
        http.Error(w, "Could not load HTML page: "+err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, week)
}

// Create a new todo for a specific day
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        title := r.FormValue("title")
        dateStr := r.FormValue("date")
        date, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            http.Error(w, "Invalid date format", http.StatusBadRequest)
            return
        }
        currentID++
        todo := Todo{ID: currentID, Title: title, Done: false, Day: date}
        todos = append(todos, todo)
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Toggle the Done status of a todo
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

// Delete a todo by ID
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
    http.HandleFunc("/", homePageHandler)
    http.HandleFunc("/todo/create", createTodoHandler)
    http.HandleFunc("/todo/toggle", toggleDoneHandler)
    http.HandleFunc("/todo/delete", deleteTodoHandler)

    fmt.Println("Server running at http://localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Error starting server:", err)
    }
}
