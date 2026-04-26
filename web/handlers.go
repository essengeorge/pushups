package web

import (
	"encoding/json"
	"net/http"
	"punkpushups/auth"
	"punkpushups/db"
	"punkpushups/domain"
	"strconv"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err := db.CreateUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, "User already exists or DB error", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	userID, err := db.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	tokenStr, err := auth.GenerateToken(userID)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}

func AddPushupsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err := db.AddPushups(userID, req.Count); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	daysStr := r.URL.Query().Get("days")
	days, _ := strconv.Atoi(daysStr)
	stats, err := db.GetStatsByUserAndFriends(userID, days)
	if err != nil {
		http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func SendFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req domain.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	recipientID, err := db.UserID(req.TargetUsername)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if err := db.SendFriendRequest(userID, recipientID); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func AcceptFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req domain.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	recipientID, err := db.UserID(req.TargetUsername)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if err := db.AcceptFriendRequest(userID, recipientID); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func RejectFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req domain.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	recipientID, err := db.UserID(req.TargetUsername)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if err := db.RejectFriendRequest(userID, recipientID); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func BlockUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req domain.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	recipientID, err := db.UserID(req.TargetUsername)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if err := db.BlockUser(userID, recipientID); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func UnblockUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req domain.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	recipientID, err := db.UserID(req.TargetUsername)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if err := db.UnblockUser(userID, recipientID); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func RemoveFriendHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req domain.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	recipientID, err := db.UserID(req.TargetUsername)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if err := db.RemoveFriend(userID, recipientID); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func GetIncomingRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	requests, err := db.GetIncomingRequests(userID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

func GetOutgoingRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	requests, err := db.GetOutgoingRequests(userID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

func GetBlockedUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	users, err := db.GetBlockedUsers(userID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	friends, err := db.GetFriends(userID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

/*
func IsFriendHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req domain.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	recipientID, err := db.UserID(req.TargetUsername)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	flag, err := db.IsFriend(userID, recipientID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	flagStr := "false"
	if flag {
		flagStr = "true"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"flag": flagStr, "status": "ok")
}

func IsBlockedHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req domain.FriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	recipientID, err := db.UserID(req.TargetUsername)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	flag, err := db.IsBlocked(userID, recipientID)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	flagStr := "false"
	if flag {
		flagStr = "true"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"flag": flagStr, "status": "ok")
}
*/
