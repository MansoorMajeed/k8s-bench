package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// Initialize SQLite database
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT,
		completed BOOLEAN
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

// CPU-intensive task (simulate JSON marshaling)
func heavyJSONProcessing(w http.ResponseWriter, r *http.Request) {
	data := make([]Task, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = Task{
			ID:          i + 1,
			Title:       "Task " + strconv.Itoa(i+1),
			Description: "This is task number " + strconv.Itoa(i+1),
			Completed:   i%2 == 0,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Create task (I/O Bound)
func createTask(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	description := r.URL.Query().Get("desc")

	_, err := db.Exec("INSERT INTO tasks (title, description, completed) VALUES (?, ?, ?)", title, description, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"status": "Task created successfully"}`))
}

// Get tasks (I/O Bound)
func getTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, description, completed FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Title, &t.Description, &t.Completed)
		tasks = append(tasks, t)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Update task (Mixed CPU + I/O Bound)
func updateTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	completed := r.URL.Query().Get("completed")

	_, err := db.Exec("UPDATE tasks SET completed = ? WHERE id = ?", completed == "true", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"status": "Task updated successfully"}`))
}

func cpuIntensiveTask(n int) int {
	if n <= 1 {
		return n
	}
	return cpuIntensiveTask(n-1) + cpuIntensiveTask(n-2)
}

func handleCPU(w http.ResponseWriter, r *http.Request) {
	n, _ := strconv.Atoi(r.URL.Query().Get("n"))
	result := cpuIntensiveTask(n)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"result": result})
}

func main() {
	initDB()
	log.Printf("GOMAXPROCS is set to: %d\n", runtime.GOMAXPROCS(0))

	http.HandleFunc("/task/create", createTask)
	http.HandleFunc("/task/list", getTasks)
	http.HandleFunc("/task/update", updateTask)
	http.HandleFunc("/heavy-json", heavyJSONProcessing)
	http.HandleFunc("/cpu", handleCPU) // pass n as a query parameter

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	log.Println("Server started on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
