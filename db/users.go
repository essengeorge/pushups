package db

import "golang.org/x/crypto/bcrypt"

func CreateUser(username, password string) error {
	var err error
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hash))
	return err
}

func Authenticate(username, password string) (int, error) {
	var id int
	var hash string
	err := DB.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", username).Scan(&id, &hash)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return 0, err
	}
	return id, nil
}
