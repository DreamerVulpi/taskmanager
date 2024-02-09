package db

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Task struct {
	Id          int
	Title       string
	Description string
	Done        bool
}

func CreateUser(conn *sql.DB, user User) error {
	slog.Info(user.Username)
	slog.Info(user.Password_hash)
	// TODO: Create method hashing password using JWT-TOKEN
	_, err := conn.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", user.Username, user.Password_hash)
	return err
}

func CreateTask(conn *sql.DB, task Task) error {
	slog.Info(fmt.Sprintf("$1, $2, $3, $4", task.Title, task.Description, task.Done))
	_, err := conn.Exec("INSERT INTO tasks (title, description) VALUES ($1, $2)", task.Title, task.Description)
	return err
}

func GetTasks(conn *sql.DB, list_id int, user_id int) ([]Task, error) {
	rows, err := conn.Query("SELECT task.id, task.title, task.description, done FROM tasks task INNER JOIN lists_tasks list_task ON list_task.task_id = task.id INNER JOIN lists_users lu ON lu.list_id = list_task.list_id WHERE list_task.list_id = $1 AND lu.user_id = $2", list_id, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tsks := []Task{}

	for rows.Next() {
		note := Task{}
		err := rows.Scan(&note.Id, &note.Title, &note.Description, &note.Done)
		if err != nil {
			slog.Warn("Error! Can't write to structure :C")
			continue
		}
		slog.Info("Successful write to structure!")
		tsks = append(tsks, note)
	}
	return tsks, err
}

func DeleteTask(conn *sql.DB, task_id int, user_id int) error {
	_, err := conn.Exec("DELETE FROM tasks task USING lists_tasks lt, lists_users lu WHERE task.id = lt.task_id AND lt.list_id = lu.list_id AND lu.user_id = $1 AND task.id = $2", user_id, task_id)
	return err
}

func UpdateTask(conn *sql.DB, task Task, user_id int) error {
	// TODO: DO SELECTABLE UPDATE (TITLE || DESCRIPTION || DONE)
	_, err := conn.Exec("UPDATE tasks task SET title='New note', description='Just something', done=true FROM lists_tasks lt, lists_users lu WHERE task.id = lt.task_id AND lt.list_id = lu.list_id AND lu.user_id = $1 AND task.id = $2", user_id, task.Id)
	return err
}
