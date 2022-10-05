package service_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testcode/test3/service"
	"testing"
)

const (
	adminUserToken   = "admin_test_token"
	regularUserToken = "regular_user_token"
)

type MockNotifier struct{}

func (n MockNotifier) Notify(text string, userId int) {
	fmt.Println("mock notified")
}

func suite(t *testing.T) (s service.Service, tearDownCh chan struct{}) {
	db, err := sql.Open("sqlite3", "./todo_test.db")
	if err != nil {
		t.Errorf(err.Error())
	}

	service.PrepareDb(db)

	if _, err := db.Exec(`DELETE FROM users; 
						  DELETE FROM todos;
						  INSERT OR IGNORE INTO users (id, is_admin, user_token) values (1, true, ?);
						  INSERT OR IGNORE INTO users (id, is_admin, user_token) values (2, false, ?);`, adminUserToken, regularUserToken); err != nil {
		t.Fatalf(err.Error())
		panic(err)
	}

	s = service.Service{Db: db, Notifier: MockNotifier{}}

	tearDownCh = make(chan struct{})

	go func() {
		<-tearDownCh
		db.Close()
		tearDownCh <- struct{}{}
	}()

	return
}

func TestCreateUser(t *testing.T) {
	service, tearDownCh := suite(t)

	req := httptest.NewRequest(http.MethodPost, "/users/?user_token="+adminUserToken, nil)
	w := httptest.NewRecorder()

	service.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if len(data) == 0 {
		t.Errorf("no data in response")
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("response status expected to be 200, got %v", res.StatusCode)
	}

	tearDownCh <- struct{}{}
	<-tearDownCh
}

func TestCreateTODO(t *testing.T) {
	service, tearDownCh := suite(t)

	testText := "test example text"

	postBodyBytes, _ := json.Marshal(map[string]string{"text": testText})

	req := httptest.NewRequest(http.MethodPost, "/todos/?user_token="+regularUserToken, bytes.NewBuffer(postBodyBytes))
	w := httptest.NewRecorder()

	service.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	responseDict := make(map[string]string)
	json.Unmarshal(data, &responseDict)

	if responseDict["text"] != testText {
		t.Errorf("wrong data in response, expected '%v', got %v", testText, responseDict["text"])
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("response status expected to be 200, got %v", res.StatusCode)
	}

	tearDownCh <- struct{}{}
	<-tearDownCh
}

func TestDeleteUser(t *testing.T) {
	service, tearDownCh := suite(t)

	req := httptest.NewRequest(http.MethodDelete, "/users/2?user_token="+regularUserToken, nil)
	w := httptest.NewRecorder()

	service.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if len(data) == 0 {
		t.Errorf("no data in response")
	}

	if res.StatusCode != http.StatusForbidden {
		t.Errorf("response status expected to be Forbidden, got %v", res.StatusCode)
	}

	tearDownCh <- struct{}{}
	<-tearDownCh
}
