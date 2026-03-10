package repository

import (
	"database/sql"
	"fmt"
	"tasks-rest-api/config"
	"tasks-rest-api/internal/model"

	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgres(config config.Config) (*PostgresRepo, error) {
	connstr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.DBUser, config.DBPassword, config.DBName, config.DBHost, config.DBPort)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresRepo{DB: db}, nil
}

func (r *PostgresRepo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, email, password, role
		FROM users
		WHERE email = $1
	`
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	return &user, err
}

func (r *PostgresRepo) CreateUser(user *model.User) error {
	query := `
		INSERT INTO users (email, password, role) 
		VALUES ($1, $2, $3)
	`
	_, err := r.DB.Exec(query, user.Email, user.Password, user.Role)
	return err
}

func (r *PostgresRepo) CreateTask(task *model.Task) error {
	query := `
		INSERT INTO tasks (title, user_id, completed)
		VALUES ($1, $2, $3)
	`
	_, err := r.DB.Exec(query, task.Title, task.UserID, task.Completed)
	return err
}

func (r *PostgresRepo) GetTasksByUserId(userID int) ([]model.Task, error) {
	rows, err := r.DB.Query("SELECT id, title, completed FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	var tasks []model.Task

	for rows.Next() {
		var task model.Task
		task.UserID = userID
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
