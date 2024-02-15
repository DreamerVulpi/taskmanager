package db

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type List struct {
	Id          int
	Title       string
	Description string
}

func CreateList(conn *sql.DB, user_id int, list List) (int, error) {
	slog.Info(fmt.Sprintf("$1, $2, $3, $4", list.Id, list.Title, list.Description))
	row, err := conn.Query("INSERT INTO lists (title, description) VALUES ($1, $2) RETURNING id", list.Title, list.Description)
	var list_id int
	if err := row.Scan(&list_id); err != nil {
		return 0, err
	}
	_, err = conn.Exec("INSERT INTO lists_users (user_id, list_id) VALUES ($1, $2)", user_id, list_id)
	if err != nil {
		return 0, err
	}
	return list_id, err
}

func GetLists(conn *sql.DB, user_id int) ([]List, error) {
	rows, err := conn.Query("SELECT list.id, list.title, list.description FROM lists list INNER JOIN lists_users lu ON list.id = lu.list_id WHERE lu.user_id = $1", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lists := []List{}

	for rows.Next() {
		note := List{}
		err := rows.Scan(&note.Id, &note.Title, &note.Description)
		if err != nil {
			slog.Warn("Error! Can't write to structure :C")
			continue
		}
		slog.Info("Successful write to structure!")
		lists = append(lists, note)
	}
	return lists, err
}

func GetList(conn *sql.DB, user_id int, list_id int) (List, error) {
	row := conn.QueryRow("SELECT list.id, list.title, list.description FROM lists list INNER JOIN lists_users lu ON list.id = lu.list_id WHERE lu.user_id = $1 AND lu.list_id = $2", user_id, list_id)
	list := List{}
	err := row.Scan(&list.Id, &list.Title, &list.Description)
	if err != nil {
		return List{}, err
	}
	return list, err
}

func DeleteList(conn *sql.DB, user_id int, task_id int) error {
	_, err := conn.Exec("DELETE FROM lists list USING lists_users lu WHERE list.id = lu.list_id AND lu.user_id = $1 AND lu.list_id = $2", user_id, task_id)
	return err
}

func UpdateList(conn *sql.DB, user_id int, list List) error {
	// TODO: DO SELECTABLE UPDATE (TITLE || DESCRIPTION)
	_, err := conn.Exec("UPDATE lists list SET title='New Here!', description='Just do it!' FROM lists_users lu, users usr WHERE list.id = lu.list_id AND usr.id = $1 AND lu.list_id = $2;", user_id, list.Id)
	return err
}
