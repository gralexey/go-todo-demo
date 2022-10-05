package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
)

type User struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}

func (t User) ToJSON() []byte {
	bytes, _ := json.Marshal(t)
	return bytes
}

func (s *Service) CreateUser(userToken string) (*User, error) {
	var userId int
	var isAdmin bool
	err := s.Db.QueryRow("SELECT id, is_admin FROM users WHERE user_token = ?", userToken).Scan(&userId, &isAdmin)
	if err != nil {
		return nil, NewErrorFromDBError(err)
	}

	if userId == 0 {
		return nil, NewErrorNotFound(nil)
	} else if !isAdmin {
		return nil, NewErrorAccessDenied()
	}

	token := generateToken()
	result, err := s.Db.Exec("INSERT INTO users (is_admin, user_token) values (false, ?)", token)
	if err != nil {
		return nil, NewErrorFromDBError(err)
	}

	newId, _ := result.LastInsertId()

	user := User{Id: int(newId), Token: token}

	return &user, nil
}

func (s *Service) DeleteUser(userToken string, targetUserId int) error {
	var userId int
	var isAdmin bool
	s.Db.QueryRow("SELECT id, is_admin FROM users WHERE user_token = ?", userToken).Scan(&userId, &isAdmin)

	if userId == 0 {
		return NewErrorNotFound(nil)
	} else if !isAdmin {
		return NewErrorAccessDenied()
	}

	if _, err := s.Db.Exec("DELETE FROM users WHERE id = ?", targetUserId); err != nil {
		return NewErrorFromDBError(err)
	}

	return nil
}

func generateToken() string {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
