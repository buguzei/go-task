package repo

import (
	"database/sql"
	"fmt"
	"github.com/buguzei/go-task/internal/config"
	"github.com/buguzei/go-task/internal/models"
	_ "github.com/lib/pq"
	"log"
)

type Postgres struct {
	DB *sql.DB
}

// NewPostgres is a constructor for Postgres
func NewPostgres(cfg config.DBConf) Postgres {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return Postgres{
		DB: db,
	}
}

func (p Postgres) GetOrderProducts(id int) ([]models.Product, error) {
	var products []models.Product

	rows, err := p.DB.Query("SELECT orders.id, order_products.product_id, order_products.amount, products.name FROM orders JOIN order_products ON orders.id = order_products.id and orders.id = ($1) JOIN products ON order_products.product_id = products.id ", id)
	if err != nil {
		return nil, fmt.Errorf("postgres query error: %w", err)
	}

	for rows.Next() {
		var prod models.Product

		err = rows.Scan(&prod.OrderID, &prod.ID, &prod.Amount, &prod.Name)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		products = append(products, prod)
	}

	for i := range products {
		row := p.DB.QueryRow("SELECT racks.name FROM products JOIN main_racks ON products.id = main_racks.product_id AND products.id = ($1) JOIN racks ON main_racks.rack_id = racks.id", products[i].ID)
		if err != nil {
			return nil, fmt.Errorf("postgres query error: %w", err)
		}

		err = row.Scan(&products[i].MainRack)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		rows, err := p.DB.Query("SELECT racks.name FROM products JOIN secondary_racks ON products.id = secondary_racks.product_id AND products.id = ($1) JOIN racks ON secondary_racks.rack_id = racks.id", products[i].ID)
		if err != nil {
			return nil, fmt.Errorf("postgres query error: %w", err)
		}

		for rows.Next() {
			var rack string
			err = rows.Scan(&rack)
			if err != nil {
				return nil, fmt.Errorf("scan error: %w", err)
			}

			products[i].SecondaryRacks = append(products[i].SecondaryRacks, rack)
		}
	}

	return products, nil
}
