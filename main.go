package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "path"
    "sync"
)

var (
    names      = []string{}
    namesMutex sync.Mutex
)

//go:embed static/*
var staticFiles embed.FS

func main() {
    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/ping", pingHandler)
    http.HandleFunc("/add", handleAdd)
    http.HandleFunc("/delete", handleDelete)
    http.HandleFunc("/names", handleNames)
    //http.Handle("/static/", http.FileServer(http.Dir(os.Getenv("KO_DATA_PATH"))))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))

    log.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Could not start server: %s\n", err)
    }
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Pong!")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
    log.Println("Serving index.html")
    indexPath := path.Join("static", "index.html")
    content, err := staticFiles.ReadFile(indexPath)
    if err != nil {
        log.Printf("Error reading index.html: %v", err)
        http.Error(w, "Page not found", http.StatusNotFound)
        return
    }
    w.Write(content)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling add request")
    if r.Method == "POST" {
        log.Println("Received POST request to add name")
        log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
        if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
            log.Println("Error parsing multipart form:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        name := r.FormValue("name")
        log.Printf("Received name: %s", name)
        if name != "" {
            namesMutex.Lock()
            names = append(names, name)
            namesMutex.Unlock()
            log.Println("Added name:", name)
            log.Println("Current names:", names)
        } else {
            log.Println("Name is empty, not adding")
        }
        jsonResponse(w, names)
    } else {
        log.Println("Invalid request method for add")
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling delete request")
    if r.Method == "POST" {
        log.Println("Received POST request to delete name")
        if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
            log.Println("Error parsing multipart form:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        name := r.FormValue("name")
        log.Println("Received name to delete:", name)
        namesMutex.Lock()
        for i, n := range names {
            if n == name {
                names = append(names[:i], names[i+1:]...)
                break
            }
        }
        namesMutex.Unlock()
        log.Println("Deleted name:", name)
        log.Println("Current names:", names)
        jsonResponse(w, names)
    } else {
        log.Println("Invalid request method for delete")
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func handleNames(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling names request")
    jsonResponse(w, names)
    log.Println("Current names:", names)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    if data == nil {
        log.Println("jsonResponse data is nil")
    } else {
        log.Println("jsonResponse data:", data)
    }
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Println("Error encoding JSON:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
