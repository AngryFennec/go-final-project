package db

import "database/sql"

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
		return []*Task{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []*Task{}, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}

	return tasks, nil
}
