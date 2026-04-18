package main

import (
        "log"
        "net/http"
        "os"
        "punkpushups/db"
        "punkpushups/web"
)

func corsMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Access-Control-Allow-Origin", "*")
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
                w.Header().Set("Access-Control-Max-Age", "86400")
                if r.Method == http.MethodOptions {
                        w.WriteHeader(http.StatusOK)
                        return
                }
                next.ServeHTTP(w, r)
        })
}

func main() {
        log.Print("Token length: ")
        log.Println(len(os.Getenv("PUSHUPS_JWT_KEY")))
        if err := db.InitDB("./pushups.db"); err != nil {
                log.Fatal("DB init failed:", err)
        }
        log.Println("Database initialized...")
        mux := http.NewServeMux()
        fileServer := http.FileServer(http.Dir("./static"))
        mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
        mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
                http.ServeFile(w, r, "./static/index.html")
        })
        mux.HandleFunc("POST /register", web.RegisterHandler)
        mux.HandleFunc("POST /login", web.LoginHandler)
        mux.HandleFunc("GET /stats", web.GetStatsHandler)
        pushupsHandler := http.HandlerFunc(web.AddPushupsHandler)
        mux.Handle("POST /pushups", web.AuthMiddleware(pushupsHandler))
        handler := corsMiddleware(mux)
        log.Println("Server started on :8080")
        if err := http.ListenAndServe(":8080", handler); err != nil {
                log.Fatal(err)
        }
}
