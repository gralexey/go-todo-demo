package service

import (
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type TodoRecord struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

func (t TodoRecord) ToJSON() []byte {
	bytes, _ := json.Marshal(t)
	return bytes
}

func (s *Service) CreateTODO(userToken string, text string) (*TodoRecord, error) {
	var userId int
	s.Db.QueryRow("SELECT id FROM users WHERE user_token = ?", userToken).Scan(&userId)

	if userId == 0 {
		return nil, NewErrorNotFound(nil)
	}

	result, err := s.Db.Exec("INSERT INTO todos (content, user_id) values (?, ?)", text, userId)
	if err != nil {
		return nil, NewErrorFromDBError(err)
	}

	newId, _ := result.LastInsertId()

	todo := &TodoRecord{Id: int(newId), Text: text}

	s.Notifier.Notify(todo.Text, userId)

	return todo, nil
}

func (s *Service) View(userToken string, id int) (*TodoRecord, error) {
	todo := TodoRecord{}
	err := s.Db.QueryRow(`WITH tmp (isAdmin, userId) AS (SELECT is_admin, id from users WHERE user_token = ?) 
							SELECT id, content FROM todos, tmp 
								WHERE (user_id = tmp.userId OR tmp.isAdmin = true) AND id = ?;`, userToken, id).Scan(&todo.Id, &todo.Text)
	if err != nil {
		return nil, NewErrorFromDBError(err)
	}

	return &todo, nil
}

func (s *Service) ViewMany(userToken string) ([]TodoRecord, error) {
	rows, err := s.Db.Query(`WITH tmp (isAdmin, userId) AS (SELECT is_admin, id from users WHERE user_token = ?) 
							SELECT id, content FROM todos, tmp 
								WHERE (user_id = tmp.userId OR tmp.isAdmin = true);`, userToken)

	if err != nil {
		return nil, NewErrorFromDBError(err)
	}
	defer rows.Close()

	todos := make([]TodoRecord, 0)
	for rows.Next() {
		todo := TodoRecord{}
		rows.Scan(&todo.Id, &todo.Text)
		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *Service) Update(userToken string, id int, text string) (*TodoRecord, error) {
	var userId int
	var isAdmin int

	s.Db.QueryRow("SELECT id, is_admin FROM users WHERE user_token = ?", userToken).Scan(&userId, &isAdmin)

	if userId == 0 {
		return nil, NewErrorNotFound(nil)
	}

	res, err := s.Db.Exec("UPDATE todos SET content = ? WHERE (user_id = ? OR ?) AND id = ?", text, userId, isAdmin, id)
	if err != nil {
		return nil, NewErrorFromDBError(err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, NewErrorFromDBError(fmt.Errorf("no data to update"))
	}

	return &TodoRecord{Id: id, Text: text}, nil
}

func (s *Service) Delete(userToken string, id int) error {
	var userId int
	var isAdmin int

	selectErr := s.Db.QueryRow("SELECT id, is_admin FROM users WHERE user_token = ?", userToken).Scan(&userId, &isAdmin)

	if selectErr != nil {
		return NewErrorFromDBError(selectErr)
	}

	if userId == 0 {
		return NewErrorNotFound(nil)
	}

	_, err := s.Db.Exec("DELETE FROM todos WHERE id = ? AND (user_id = ? OR ?)", id, userId, isAdmin)
	if err != nil {
		return NewErrorFromDBError(err)
	}
	return nil
}
