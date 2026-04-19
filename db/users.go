package db

import "golang.org/x/crypto/bcrypt"

const (
	codeBanned   = -1
	codePending  = 0
	codeApproved = 1
)

func CreateUser(username, password string) error {
	var err error
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hash))
	return err
}

func BanUser(username string) error {
	_, err := DB.Exec("UPDATE users SET is_approved = ? WHERE username = ?", codeBanned, username)
	return err
}

func ApproveUser(username string) error {
	_, err := DB.Exec("UPDATE users SET is_approved = ? WHERE username = ?", codeApproved, username)
	return err
}

func UserRole(username string) (string, error) {
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

func UserRoleByID(userID int) (string, error) {
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}


func Authenticate(username, password string) (int, error) {
	var id int
	var hash string
	err := DB.QueryRow("SELECT id, password_hash FROM users WHERE username = ? AND is_approved = ?", username, codeApproved).Scan(&id, &hash)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetPendingUsers() ([]string, error) {
	rows, err := DB.Query("SELECT username FROM users WHERE is_approved = ?", codePending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	usernames := []string{}
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}
	return usernames, nil
}
