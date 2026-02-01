package repositories

import (
	"context"
	"errors"
	"kasir-api/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) GetAllCategories() ([]models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `SELECT id, name, description FROM categories ORDER BY id`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (r *CategoryRepository) CreateCategory(c *models.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	return r.pool.QueryRow(ctx, q, c.Name, c.Description).Scan(&c.ID)
}

func (r *CategoryRepository) GetCategoryByID(id int) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `SELECT id, name, description FROM categories WHERE id = $1`
	var c models.Category
	if err := r.pool.QueryRow(ctx, q, id).Scan(&c.ID, &c.Name, &c.Description); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) UpdateCategory(c *models.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `UPDATE categories SET name = $1, description = $2 WHERE id = $3`
	ct, err := r.pool.Exec(ctx, q, c.Name, c.Description, c.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *CategoryRepository) DeleteCategory(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const q = `DELETE FROM categories WHERE id = $1`
	ct, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("category not found")
	}
	return nil
}
