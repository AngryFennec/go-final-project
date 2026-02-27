package db

import (
	"database/sql"
	"errors"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(t *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(limit int) ([]*Task, error) {
	tasks := make([]*Task, 0)
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`
	rows, err := DB.Query(query, sql.Named("limit", limit))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task *Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id;`

	err := DB.QueryRow(query, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, repeat = :repeat, comment = :comment WHERE id = :id`
	res, err := DB.Exec(query, sql.Named("id", task.ID), sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("repeat", task.Repeat), sql.Named("comment", task.Comment))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id;`
	res, err := DB.Exec(query, sql.Named("id", id))

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("task not found")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id;`
	res, err := DB.Exec(query, sql.Named("id", id), sql.Named("date", next))

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("task not found")
	}

	return nil
}
