package repositories

import (
	"context"
	"kasir-api/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	pool *pgxpool.Pool
}

func NewReportRepository(pool *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{pool: pool}
}

func (r *ReportRepository) GetTodayReport() (*models.TodayReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var totalRevenue int
	if err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total_amount), 0)
		FROM transactions t
		JOIN transaction_details d ON d.transaction_id = t.id
		WHERE t.created_at::date = CURRENT_DATE
	`).Scan(&totalRevenue); err != nil {
		return nil, err
	}

	var totalTransaction int
	if err := r.pool.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM transactions t
        WHERE t.created_at::date = CURRENT_DATE
    `).Scan(&totalTransaction); err != nil {
		return nil, err
	}

	var name string
	var qty int
	err := r.pool.QueryRow(ctx, `
        SELECT p.name, SUM(d.quantity) AS qty_sold
        FROM transactions t
        JOIN transaction_details d ON d.transaction_id = t.id
        JOIN products p ON p.id = d.product_id
        WHERE t.created_at::date = CURRENT_DATE
        GROUP BY p.name
        ORDER BY qty_sold DESC
        LIMIT 1
    `).Scan(&name, &qty)
	if err == pgx.ErrNoRows {
		name, qty = "", 0
	} else if err != nil {
		return nil, err
	}

	return &models.TodayReport{
		TotalRevenue:      totalRevenue,
		TotalTransactions: totalTransaction,
		BestsellingProducts: []models.BestsellingProduct{
			{Name: name, QtySold: qty},
		},
	}, nil
}
