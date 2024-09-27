// Prueba001 uSZ%c@YRKE@G{1.N
// mycrud-golan:us-central1:my-postgres-crud-golan

package db

import (
	"context"
	"fmt"
	"log"
	url "net/url"
	"os"

	"My-CRUD-Golang/internal/domain"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository() *PostgresRepository {
	dbUser := getenv("DB_USER", "default_user")
	dbPassword := getenv("DB_PASSWORD", "default_password")
	dbName := getenv("DB_NAME", "crud_golang_items")
	instanceConnectionName := getenv("INSTANCE_CONNECTION_NAME", "project-id:region:instance-id")

	encodedUser := url.QueryEscape(dbUser)
	encodedPassword := url.QueryEscape(dbPassword)

	dsn := fmt.Sprintf("postgres://%s:%s@/%s?host=/cloudsql/%s",
		encodedUser, encodedPassword, dbName, instanceConnectionName)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}

	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return &PostgresRepository{
		db: dbpool,
	}
}

func (r *PostgresRepository) GetAll() ([]domain.Item, error) {
	rows, err := r.db.Query(context.Background(), "SELECT id, name, description FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		err := rows.Scan(&item.ID, &item.Name, &item.Price)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *PostgresRepository) GetByID(id string) (*domain.Item, error) {
	var item domain.Item
	err := r.db.QueryRow(context.Background(), "SELECT id, name, description FROM items WHERE id=$1", id).
		Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		return &item, err
	}
	return &item, nil
}

func (r *PostgresRepository) Create(item domain.Item) error {
	_, err := r.db.Exec(context.Background(), "INSERT INTO items (name, description) VALUES ($1, $2)",
		item.Name, item.Price)
	return err
}

func (r *PostgresRepository) Update(item domain.Item) error {
	_, err := r.db.Exec(context.Background(), "UPDATE items SET name=$1, description=$2 WHERE id=$3",
		item.Name, item.Price, item.ID)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM items WHERE id=$1", id)
	return err
}

func getenv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
