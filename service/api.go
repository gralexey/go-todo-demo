package service

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	urlComps := strings.Split(url, "/")
	id := -1
	if len(urlComps) > 2 {
		id, _ = strconv.Atoi(urlComps[2])
	}

	userToken := r.URL.Query().Get("user_token")

	if len(userToken) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no user token provided"))
		return
	}

	body, _ := io.ReadAll(r.Body)

	defer r.Body.Close()
	bodyJson := make(map[string]string)
	json.Unmarshal(body, &bodyJson)
	text := bodyJson["text"]

	var serviceError *Error

	if strings.HasPrefix(url, "/todos") {
		switch r.Method {
		case "GET":
			if id == 0 {
				todos, err := s.ViewMany(userToken)
				if errors.As(err, &serviceError) {
					http.Error(w, err.Error(), serviceError.suggestedCode)
				} else {
					json, _ := json.Marshal(todos)
					w.Write(json)
				}
			} else {
				todo, err := s.View(userToken, id)
				if errors.As(err, &serviceError) {
					http.Error(w, err.Error(), serviceError.suggestedCode)
				} else {
					w.Write(todo.ToJSON())
				}
			}

		case "POST":
			todo, err := s.CreateTODO(userToken, text)
			if errors.As(err, &serviceError) {
				http.Error(w, err.Error(), serviceError.suggestedCode)
			} else {
				w.Write(todo.ToJSON())
			}

		case "DELETE":
			if err := s.Delete(userToken, id); errors.As(err, &serviceError) {
				http.Error(w, serviceError.text, serviceError.suggestedCode)
			}

		case "PATCH":
			todo, err := s.Update(userToken, id, text)
			if errors.As(err, &serviceError) {
				http.Error(w, err.Error(), serviceError.suggestedCode)
			} else {
				w.Write(todo.ToJSON())
			}

		default:
			log.Println("unexpected method")
		}

	} else if strings.HasPrefix(url, "/users") {
		switch r.Method {
		case "POST":
			user, err := s.CreateUser(userToken)
			if errors.As(err, &serviceError) {
				http.Error(w, err.Error(), serviceError.suggestedCode)
			} else {
				w.Write(user.ToJSON())
			}

		case "DELETE":
			if errors.As(s.DeleteUser(userToken, id), &serviceError) {
				http.Error(w, serviceError.Error(), serviceError.suggestedCode)
			}

		default:
			log.Println("unexpected method")
		}
	}
}
