package repositories

import (
	"context"
	"errors"
	"kasir-api/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `SELECT id, name, price, stock FROM products`
	rows, err := repo.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (repo *ProductRepository) Create(product *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`
	return repo.pool.QueryRow(ctx, query, product.Name, product.Price, product.Stock).Scan(&product.ID)
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `SELECT id, name, price, stock FROM products WHERE id = $1`
	var p models.Product
	err := repo.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		// pgx tidak mengembalikan sql.ErrNoRows, gunakan errors.New untuk pesan user-friendly
		return nil, errors.New("product not found")
	}
	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4`
	ct, err := repo.pool.Exec(ctx, query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `DELETE FROM products WHERE id = $1`
	ct, err := repo.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}
