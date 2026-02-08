package repositories

import (
	"context"
	"database/sql"
	"errors"
	"kasir-api/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

	const query = `
		INSERT INTO products (name, price, stock, category_id) 
		VALUES ($1, $2, $3, $4) RETURNING id`
	return repo.pool.QueryRow(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, COALESCE(c.name, '') AS category_name
        FROM products p
        LEFT JOIN categories c ON c.id = p.category_id
        WHERE p.id = $1
		`
	var (
		p       models.Product
		catID   sql.NullInt32
		catName string
	)

	err := repo.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID, &catName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	if catID.Valid {
		v := int(catID.Int32)
		p.CategoryID = &v
	} else {
		p.CategoryID = nil
	}
	p.CategoryName = catName

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// category_id nullable
	var cat pgtype.Int4
	if product.CategoryID != nil {
		cat = pgtype.Int4{Int32: int32(*product.CategoryID), Valid: true}
	} else {
		cat = pgtype.Int4{Valid: false} // akan ditulis sebagai NULL
	}

	const query = `UPDATE products 
				   SET name = $1, price = $2, stock = $3, category_id = $4 
				   WHERE id = $5`
	ct, err := repo.pool.Exec(ctx, query, product.Name, product.Price, product.Stock, cat, product.ID)
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
