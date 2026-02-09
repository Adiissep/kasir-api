package repositories

import (
	"context"
	"errors"
	"kasir-api/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{pool: pool}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := repo.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	totalAmount := 0
	details := make([]models.TransactionDetail, 0, len(items))

	for _, item := range items {
		var productName string
		var productPrice int
		var stock int

		// get data product
		err = tx.QueryRow(ctx, `
            SELECT name, price, stock
            FROM products
            WHERE id = $1
        `, item.ProductID).Scan(&productName, &productPrice, &stock)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.New("product id not found")
			}
			return nil, err
		}
		// validation stoct
		if stock < item.Quantity {
			return nil, errors.New("insufficient stock")
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// Reduce stock
		_, err = tx.Exec(ctx, `
            UPDATE products
            SET stock = stock - $1
            WHERE id = $2
        `, item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Insert transaction
	var transactionID int
	err = tx.QueryRow(ctx, `
        INSERT INTO transactions (total_amount)
        VALUES ($1)
        RETURNING id
    `, totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Insert transaction detail
	for i := range details {
		details[i].TransactionID = transactionID
		var detailID int
		err = tx.QueryRow(ctx, `
            INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal)
			VALUES ($1, $2, $3, $4)
			RETURNING id
        `, transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
