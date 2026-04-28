package main

import (
	"log"
	"net/http"
	"os"
	"punkpushups/db"
	"punkpushups/web"
)

func main() {
	log.Println("JWT key length:")
	jwtKeyLength := len(os.Getenv("PUSHUPS_JWT_KEY"))
	if jwtKeyLength < 32 {
		log.Fatal("Too weak JWT key!")
	}
	log.Println(jwtKeyLength)
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

	adminFileServer := http.FileServer(http.Dir("./static/admin"))
	mux.Handle("GET /static/admin", http.StripPrefix("/static/admin", adminFileServer))
	mux.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/admin/index.html")
	})

	mux.HandleFunc("POST /register", web.RegisterHandler)
	mux.HandleFunc("POST /login", web.LoginHandler)

	getStatsHandler := http.HandlerFunc(web.GetStatsHandler)
	mux.Handle("GET /stats", web.AuthMiddleware(web.ApprovedMiddleware(getStatsHandler)))

	pushupsHandler := http.HandlerFunc(web.AddPushupsHandler)
	mux.Handle("POST /pushups", web.AuthMiddleware(web.ApprovedMiddleware(pushupsHandler)))

	sendFriendRequestHandler := http.HandlerFunc(web.SendFriendRequestHandler)
	mux.Handle("POST /friends/send", web.AuthMiddleware(web.ApprovedMiddleware(sendFriendRequestHandler)))

	acceptFriendRequestHandler := http.HandlerFunc(web.AcceptFriendRequestHandler)
	mux.Handle("POST /friends/accept", web.AuthMiddleware(web.ApprovedMiddleware(acceptFriendRequestHandler)))

	rejectFriendRequestHandler := http.HandlerFunc(web.RejectFriendRequestHandler)
	mux.Handle("POST /friends/reject", web.AuthMiddleware(web.ApprovedMiddleware(rejectFriendRequestHandler)))

	blockUserHandler := http.HandlerFunc(web.BlockUserHandler)
	mux.Handle("POST /friends/block", web.AuthMiddleware(web.ApprovedMiddleware(blockUserHandler)))

	unblockUserHandler := http.HandlerFunc(web.UnblockUserHandler)
	mux.Handle("POST /friends/unblock", web.AuthMiddleware(web.ApprovedMiddleware(unblockUserHandler)))

	removeFriendHandler := http.HandlerFunc(web.RemoveFriendHandler)
	mux.Handle("POST /friends/remove", web.AuthMiddleware(web.ApprovedMiddleware(removeFriendHandler)))

	getIncomingRequestsHandler := http.HandlerFunc(web.GetIncomingRequestsHandler)
	mux.Handle("GET /friends/incoming", web.AuthMiddleware(web.ApprovedMiddleware(getIncomingRequestsHandler)))

	getOutgoingRequestsHandler := http.HandlerFunc(web.GetOutgoingRequestsHandler)
	mux.Handle("GET /friends/outgoing", web.AuthMiddleware(web.ApprovedMiddleware(getOutgoingRequestsHandler)))

	getBlockedUsersHandler := http.HandlerFunc(web.GetBlockedUsersHandler)
	mux.Handle("GET /friends/blocked", web.AuthMiddleware(web.ApprovedMiddleware(getBlockedUsersHandler)))

	getFriendsHandler := http.HandlerFunc(web.GetFriendsHandler)
	mux.Handle("GET /friends/list", web.AuthMiddleware(web.ApprovedMiddleware(getFriendsHandler)))

	bannedHandler := http.HandlerFunc(web.GetBannedUsersHandler)
	mux.Handle("GET /admin/banned", web.AuthMiddleware(web.ApprovedMiddleware(web.AdminMiddleware(bannedHandler))))

	pendingHandler := http.HandlerFunc(web.GetPendingUsersHandler)
	mux.Handle("GET /admin/pending", web.AuthMiddleware(web.ApprovedMiddleware(web.AdminMiddleware(pendingHandler))))

	approvedHandler := http.HandlerFunc(web.GetApprovedUsersHandler)
	mux.Handle("GET /admin/approved", web.AuthMiddleware(web.ApprovedMiddleware(web.AdminMiddleware(approvedHandler))))

	approveHandler := http.HandlerFunc(web.ApproveUserHandler)
	mux.Handle("POST /admin/approve", web.AuthMiddleware(web.ApprovedMiddleware(web.AdminMiddleware(approveHandler))))

	banHandler := http.HandlerFunc(web.BanUserHandler)
	mux.Handle("POST /admin/ban", web.AuthMiddleware(web.ApprovedMiddleware(web.AdminMiddleware(banHandler))))

	handler := web.CorsMiddleware(mux)
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
