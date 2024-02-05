package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

type CRUD interface {
	CREATE(data note)
	READ()
	UPDATE(data note)
	DELETE(id int)
}

type note struct {
	Id        int
	Objective string
	Deadline  string
}

// TODO: CRUS

// func CREATE(conn *sql.DB, data note) {
// 	slog.Info(fmt.Sprintf("$1, $2", data.Objective, data.Deadline))

// }

func READ(conn *sql.DB) []note {
	rows, err := conn.Query("SELECT * FROM task")
	if err != nil {
		slog.Warn("Error. Can't read table :C")
	}
	defer rows.Close()

	notes := []note{}

	for rows.Next() {
		note := note{}
		err := rows.Scan(&note.Id, &note.Objective, &note.Deadline)
		if err != nil {
			slog.Warn("Error! Can't write to structure :C")
			continue
		}
		slog.Info(note.Objective + " " + note.Deadline)
		notes = append(notes, note)
	}

	return notes
}

// func UPDATE(conn *sql.DB, data inputData) {

// }

// func DELETE(conn *sql.DB, id int) {

// }
