package main

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"scrum/internal/api"
)

type GetUsersResponse struct {
	Count int64
	Users []User
}
type User struct {
	Username string `json:"username,omitempty"`
}

func main() {
	http.HandleFunc("/board/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.BoardList(w, r)
	})

	http.HandleFunc("/board/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.СreateBoard(w, r)
	})

	// TODO: Создать Git репозиторий и добавить в него проект

	// TODO: Добавить редактирование доски
	// TODO: Добавить Арифметику идеального рабочего времени (estimation) последней

	// TODO: Добавить Создание карточки доски
	// TODO: Добавить Редактирование карточки доски
	// TODO: Добавить Удаление карточки доски
	// TODO: Добавить Получения списка карточек

	// TODO: Добавить Авторизацию (в каждый запрос нужно передавать login|password, проверять их в БД и возвращать наверх его id,
	// что бы передать его в запрос. Например создание доски должно требовать передачи поля creator_id)

	// TODO: Добавить Отчет по колонке

	http.HandleFunc("/board/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.DeleteBoard(w, r)
	})

	http.HandleFunc("/user/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.UsersList(w, r)
	})

	http.HandleFunc("/card/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.СreateCard(w, r)
	})

	http.HandleFunc("/card/update", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.UpdateCard(w, r)
	})

	http.HandleFunc("/card/delete", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.DeleteCard(w, r)
	})

	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		api.Report(w, r)
	})

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{'status': 'ok'}"))
	})

	addres := "0.0.0.0:8888"
	fmt.Println(fmt.Sprintf("run server: %s", addres))
	http.ListenAndServe(addres, nil)
}
