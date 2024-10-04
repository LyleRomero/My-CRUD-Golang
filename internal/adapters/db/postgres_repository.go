// Prueba001 uSZ%c@YRKE@G{1.N
// mycrud-golan:us-central1:my-postgres-crud-golan

package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"My-CRUD-Golang/internal/domain"
	"My-CRUD-Golang/internal/ports"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository crea una nueva instancia de PostgresRepository
func NewPostgresRepository(db *sql.DB) ports.ItemRepository {
	return &PostgresRepository{
		db: db,
	}
}

func ConnectWithConnectorIAMAuthN() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'service-account-name@project-id.iam'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	d, err := cloudsqlconn.NewDialer(context.Background(), cloudsqlconn.WithIAMAuthN())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}
	var opts []cloudsqlconn.DialOption
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}

	dsn := fmt.Sprintf("user=%s database=%s", dbUser, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName, opts...)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return dbPool, nil
}

// Implementación de los métodos de la interfaz ItemRepository

func (r *PostgresRepository) GetAll() ([]domain.Item, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *PostgresRepository) GetByID(id string) (*domain.Item, error) {
	var item domain.Item
	err := r.db.QueryRow("SELECT id, name, description FROM items WHERE id = $1", id).
		Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PostgresRepository) Create(item domain.Item) error {
	_, err := r.db.Exec("INSERT INTO items (id, name, description) VALUES ($1, $2, $3)", item.ID, item.Name, item.Price)
	return err
}

func (r *PostgresRepository) Update(item domain.Item) error {
	_, err := r.db.Exec("UPDATE items SET name = $1, description = $2 WHERE id = $3", item.Name, item.Price, item.ID)
	return err
}

func (r *PostgresRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM items WHERE id = $1", id)
	return err
}
