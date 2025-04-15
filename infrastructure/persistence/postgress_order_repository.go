package persistence

import (
	"ORDERING-API/domain/entities"
	"ORDERING-API/domain/repositories"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type PostgresOrderRepository struct {
	db *sql.DB
}

// NewOrderRepository initializes a new repository
func NewOrderRepository(db *sql.DB) repositories.OrderRepository {
	return &PostgresOrderRepository{db: db}
}

// SaveOrder inserts an order into the database
func (r *PostgresOrderRepository) SaveOrder(order *entities.Order) (string, error) {

	query := `INSERT INTO orders (id, notes, total) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, order.Id, order.Notes, order.Total)
	if err != nil {
		return "", fmt.Errorf("failed to save order: %v", err)
	}

	// Insert order items
	for _, item := range order.OrderItems {
		itemQuery := `INSERT INTO orderitems (id, orderid, productid,  quantity) VALUES ($1, $2, $3, $4)`
		_, err := r.db.Exec(itemQuery, item.Id, order.Id, item.ProductID, item.Quantity)
		if err != nil {
			return "", fmt.Errorf("failed to save order item: %v", err)
		}
	}

	return order.Id, nil
}

// GetOrderByID retrieves an order from the database
func (r *PostgresOrderRepository) GetOrderByID(id string) (*entities.Order, error) {
	var order entities.Order
	query := `SELECT id, notes, total FROM orders WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&order.Id, &order.Notes, &order.Total)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Retrieve order items
	itemQuery := `SELECT id, orderid, productid,  quantity FROM orderitems WHERE orderid = $1`
	rows, err := r.db.Query(itemQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entities.OrderItem
	for rows.Next() {
		var item entities.OrderItem
		if err := rows.Scan(&item.Id, &item.OrderID, &item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	order.OrderItems = items

	return &order, nil
}

func (r *PostgresOrderRepository) UpdateOrder(order *entities.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Update order
	orderQuery := `UPDATE orders SET notes = $1, total = $2 WHERE id = $3`
	_, err = tx.Exec(orderQuery, order.Notes, order.Total, order.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update order: %v", err)
	}

	// Delete existing order items
	deleteItemsQuery := `DELETE FROM orderitems WHERE orderid = $1`
	_, err = tx.Exec(deleteItemsQuery, order.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete old order items: %v", err)
	}

	// Insert new order items
	for _, item := range order.OrderItems {
		itemQuery := `INSERT INTO orderitems (id, orderid, productid, quantity) VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(itemQuery, item.Id, order.Id, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert updated order items: %v", err)
		}
	}

	return tx.Commit()
}

// DeleteOrder removes an order and its items from the database
func (r *PostgresOrderRepository) DeleteOrder(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Delete order items first
	deleteItemsQuery := `DELETE FROM orderitems WHERE orderid = $1`
	_, err = tx.Exec(deleteItemsQuery, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete order items: %v", err)
	}

	// Delete order
	orderQuery := `DELETE FROM orders WHERE id = $1`
	_, err = tx.Exec(orderQuery, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete order: %v", err)
	}

	return tx.Commit()
}

// GetOrdersWithPagination retrieves paginated orders
func (r *PostgresOrderRepository) GetOrdersWithPagination(limit, offset int) ([]entities.Order, error) {
	query := `SELECT id, notes, total FROM orders LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entities.Order
	for rows.Next() {
		var order entities.Order
		if err := rows.Scan(&order.Id, &order.Notes, &order.Total); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
