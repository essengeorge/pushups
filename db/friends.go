package db

import (
	//	"database/sql"
	"errors"
	"punkpushups/domain"
)

const (
	friendsCodeBlocked  = -1
	friendsCodePending  = 0
	friendsCodeAccepted = 1
)

func SendFriendRequest(requesterID, recipientID int) error {
	if requesterID == recipientID {
		return errors.New("cannot send friend request to yourself")
	}
	_, err := DB.Exec(
		"INSERT INTO friendships (requester_id, recipient_id, status) VALUES (?, ?, ?)",
		requesterID, recipientID, friendsCodePending,
	)
	return err
}

func AcceptFriendRequest(requesterID, recipientID int) error {
	_, err := DB.Exec(
		"UPDATE friendships SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE requester_id = ? AND recipient_id = ? AND status = ?",
		friendsCodeAccepted, recipientID, requesterID, friendsCodePending,
	)
	return err
}

func RejectFriendRequest(requesterID, recipientID int) error {
	_, err := DB.Exec(
		"DELETE FROM friendships WHERE ((requester_id = ? AND recipient_id = ?) OR (requester_id = ? AND recipient_id = ?)) AND status = ?",
		requesterID, recipientID, recipientID, requesterID, friendsCodePending,
	)
	return err
}

func BlockUser(requesterID, recipientID int) error {
	_, err := DB.Exec(`
		DELETE FROM friendships
		WHERE (requester_id = ? AND recipient_id = ?) OR (requester_id = ? AND recipient_id = ?)
	`, requesterID, recipientID, recipientID, requesterID)
	if err != nil {
		return err
	}
	_, err = DB.Exec(
		"INSERT INTO friendships (requester_id, recipient_id, status) VALUES (?, ?, ?)",
		requesterID, recipientID, friendsCodeBlocked,
	)
	return err
}

func UnblockUser(requesterID, recipientID int) error {
	_, err := DB.Exec(
		"DELETE FROM friendships WHERE requester_id = ? AND recipient_id = ? AND status = ?",
		requesterID, recipientID, friendsCodeBlocked,
	)
	return err
}

func RemoveFriend(userID, friendID int) error {
	_, err := DB.Exec(`
		DELETE FROM friendships
		WHERE ((requester_id = ? AND recipient_id = ?) OR (requester_id = ? AND recipient_id = ?)) AND (status = ?)
	`, userID, friendID, friendID, userID, friendsCodeAccepted)
	return err
}

func GetIncomingRequests(userID int) ([]domain.Friendship, error) {
	rows, err := DB.Query(`
		SELECT f.recipient_id, u.id, u.username, f.status
		FROM friendships f
		JOIN users u ON f.requester_id = u.id
		WHERE f.recipient_id = ? AND f.status = ?
		ORDER BY f.created_at DESC
	`, userID, friendsCodePending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friendships []domain.Friendship
	for rows.Next() {
		var f domain.Friendship
		if err := rows.Scan(&f.UserID, &f.FriendID, &f.Username, &f.Status); err != nil {
			return nil, err
		}
		friendships = append(friendships, f)
	}
	return friendships, nil
}

func GetOutgoingRequests(userID int) ([]domain.Friendship, error) {
	rows, err := DB.Query(`
		SELECT f.requester_id, u.id, u.username, f.status
		FROM friendships f
		JOIN users u ON f.recipient_id = u.id
		WHERE f.requester_id = ? AND f.status = ?
		ORDER BY f.created_at DESC
	`, userID, friendsCodePending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friendships []domain.Friendship
	for rows.Next() {
		var f domain.Friendship
		if err := rows.Scan(&f.UserID, &f.FriendID, &f.Username, &f.Status); err != nil {
			return nil, err
		}
		friendships = append(friendships, f)
	}
	return friendships, nil
}

func GetFriends(userID int) ([]domain.Friendship, error) {
	rows, err := DB.Query(`
		SELECT u.id, u.username
		FROM friendships f
		JOIN users u ON (
			(f.requester_id = ? AND f.recipient_id = u.id) OR
			(f.recipient_id = ? AND f.requester_id = u.id)
		)
		WHERE f.status = ?
		ORDER BY u.username
	`, userID, userID, friendsCodeAccepted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friendships []domain.Friendship
	for rows.Next() {
		var f domain.Friendship
		if err := rows.Scan(&f.FriendID, &f.Username); err != nil {
			return nil, err
		}
		f.Status = friendsCodeAccepted
		friendships = append(friendships, f)
	}
	return friendships, nil
}

func GetBlockedUsers(userID int) ([]domain.Friendship, error) {
	rows, err := DB.Query(`
		SELECT f.requester_id, u.id, u.username, f.status
		FROM friendships f
		JOIN users u ON f.recipient_id = u.id
		WHERE f.requester_id = ? AND f.status = ?
		ORDER BY f.created_at DESC
	`, userID, friendsCodeBlocked)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var friendships []domain.Friendship
	for rows.Next() {
		var f domain.Friendship
		if err := rows.Scan(&f.UserID, &f.FriendID, &f.Username, &f.Status); err != nil {
			return nil, err
		}
		friendships = append(friendships, f)
	}
	return friendships, nil
}

/*
func IsFriend(userID, otherUserID int) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) FROM friendships
		WHERE (requester_id = ? AND recipient_id = ?) OR (requester_id = ? AND recipient_id = ?)
		AND status = ?
	`, userID, otherUserID, otherUserID, userID, friendsCodeAccepted).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsBlocked(requesterID, recipientID int) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) FROM friendships
		WHERE requester_id = ? AND recipient_id = ? AND status = ?
	`, requesterID, recipientID, friendsCodeBlocked).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return count > 0, nil
}
*/
