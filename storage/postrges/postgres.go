package postrges

import (
	"Rest-shortcut/internal/config"
	"Rest-shortcut/lib/models"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(conf config.StorageConfig) (*Storage, error) {
	const op = "storage.postgres.NewStorage"
	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		conf.User, conf.Password, conf.Address, conf.Port, conf.DbName)
	db, err := sql.Open("pgx", dataSource)
	if err != nil {
		return nil, fmt.Errorf("%s Open data base: %w", op, err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s Ping data base: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (storage *Storage) AddUser(login, password string) error {
	const op = "storage.postgres.AddUser"
	query, err := storage.db.Prepare("INSERT INTO users (email, password) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s prepare query: %w", op, err)
	}

	_, err = query.Exec(login, password)
	if err != nil {
		return fmt.Errorf("%s execute query: %w", op, err)
	}
	return nil
}
func (storage *Storage) GetUser(login, password string) (models.User, error) {
	const op = "storage.postgres.GetUser"
	query, err := storage.db.Prepare("SELECT email, password FROM users WHERE email = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s prepare query: %w", op, err)
	}
	var user models.User
	err = query.QueryRow(login).Scan(&user.Login, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("%s query row: %w", op, err)
	}
	return user, nil
}
func (storage *Storage) UpdateRefreshToken(username, refreshToken string) error {
	const op = "storage.postgres.UpdateRefreshToken"
	query, err := storage.db.Prepare("UPDATE users SET refresh_token = $1 WHERE email = $2")
	if err != nil {
		return fmt.Errorf("%s prepare query: %w", op, err)
	}
	_, err = query.Exec(refreshToken, username)
	if err != nil {
		return fmt.Errorf("%s execute query: %w", op, err)
	}
	return nil
}
func (storage *Storage) SaveUrl(user, oldUrl, shortUrl string) error {
	const op = "storage.postgres.SaveUrl"
	smth, err := storage.db.Prepare("INSERT INTO url(user_id, old_url, new_url) SELECT id, $1, $2 FROM users WHERE email = $3")
	if err != nil {
		return fmt.Errorf("%s Prepare query: %w", op, err)
	}
	_, err = smth.Exec(oldUrl, shortUrl, user)
	if err != nil {
		return fmt.Errorf("%s Execute query: %w", op, err)
	}
	return nil
}

func (storage *Storage) GetUrl(user, shortUrl string) (string, error) {
	const op = "storage.postgres.GetUrl"
	query, err := storage.db.Prepare(
		"SELECT url.old_url FROM url JOIN users ON url.user_id = users.id WHERE users.email=$1 AND url.new_url=$2")

	if err != nil {
		return "", fmt.Errorf("%s Prepare query: %w", op, err)
	}
	var result string
	err = query.QueryRow(user, shortUrl).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("%s QueryRow query: %w", op, err)
	}

	return result, nil
}

func (storage *Storage) DeleteUrl(shortUrl string) error {
	const op = "storage.postgres.DeleteUrl"
	query, err := storage.db.Prepare("DELETE FROM url WHERE new_url=$1")
	if err != nil {
		return fmt.Errorf("%s Prepare query: %w", op, err)
	}
	_, err = query.Exec(shortUrl)
	if err != nil {
		return fmt.Errorf("%s Execute query: %w", op, err)
	}
	return nil
}
