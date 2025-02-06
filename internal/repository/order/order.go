package order

import (
	"context"
	"database/sql"
	"encoding/json"
	orderModel "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/order"
	orderModels "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/orderItem"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) GetUsersOrder(ctx context.Context, userId string) (*orderModel.Model, error) {
	const op = "repository.order.GetUsersOrder"

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, user_id, status, total_price FROM orders WHERE user_id = $1 and status = 'draft'")
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer stmt.Close()

	var order orderModel.Model

	row := stmt.QueryRowContext(ctx, userId)
	if row.Err() != nil {
		return nil, errors.Wrap(err, op)
	}

	err = row.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalPrice)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &order, nil
}

func (r *Repository) GetUserOrderByUserID(ctx context.Context, orderId string) (*orderModel.Model, error) {
	const op = "repository.order.GetUsersOrder"

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, user_id, status, total_price FROM orders WHERE user_id = $1")
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer stmt.Close()

	var order orderModel.Model

	row := stmt.QueryRowContext(ctx, orderId)
	if row.Err() != nil {
		return nil, errors.Wrap(err, op)
	}

	err = row.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalPrice)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &order, nil
}

func (r *Repository) CreateUserOrder(ctx context.Context, userID string) error {
	const op = "repository.order.CreateUserOrder"

	stmt, err := r.DB.PrepareContext(ctx, "INSERT INTO orders(user_id) VALUES ($1)")
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repository) CheckOrderExists(ctx context.Context, userID string) (bool, error) {
	const op = "repository.order.CheckOrderExists"

	stmt, err := r.DB.PrepareContext(ctx, "SELECT EXISTS(SELECT 1 FROM orders WHERE user_id = $1 AND status = 'draft')")
	if err != nil {
		return false, errors.Wrap(err, op)
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRowContext(ctx, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, errors.Wrap(err, op)
	}

	return exists, nil
}

func (r *Repository) AlterUserOrder(ctx context.Context, userID string) error {
	const op = "repository.order.AlterUserOrder"

	stmt, err := r.DB.PrepareContext(ctx, "UPDATE orders SET status = 'paid' WHERE user_id = ?")
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repository) AddOrderItemIntoOrder(ctx context.Context, userID string, items *[]orderModels.OrderItem) error {
	const op = "repository.order.AddOrderItemIntoOrder"

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, op+": failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	orderID, err := r.getOrCreateOrder(ctx, tx, userID)
	if err != nil {
		return err
	}

	orderStatus, err := r.checkOrderStatus(ctx, tx, orderID)
	if err != nil {
		return err
	}

	if orderStatus != "draft" {
		return errors.Wrap(errors.New("order is not in draft status"), op)
	}

	stmtBookPrice, stmtCheckExisting, stmtUpdateItem, stmtInsertItem, stmtUpdateStock, stmtUpdateTotal, err := r.prepareStatements(ctx, tx)
	if err != nil {
		return err
	}
	defer stmtBookPrice.Close()
	defer stmtCheckExisting.Close()
	defer stmtUpdateItem.Close()
	defer stmtInsertItem.Close()
	defer stmtUpdateStock.Close()
	defer stmtUpdateTotal.Close()

	for _, item := range *items {
		if err := r.processItem(ctx, tx, item, orderID, stmtBookPrice, stmtCheckExisting, stmtUpdateItem, stmtInsertItem, stmtUpdateStock, stmtUpdateTotal); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, op+": failed to commit transaction")
	}

	return nil
}

func (r *Repository) getOrCreateOrder(ctx context.Context, tx *sql.Tx, userID string) (uuid.UUID, error) {
	const op = "repository.order.getOrCreateOrder"

	exists, err := r.CheckOrderExists(ctx, userID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op+": failed to check existing order")
	}

	if exists {
		order, err := r.GetUsersOrder(ctx, userID)
		if err != nil {
			return uuid.Nil, errors.Wrap(err, op+": failed to get existing order")
		}
		return order.ID, nil
	}

	err = r.CreateUserOrder(ctx, userID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op+": failed to create new order")
	}

	order, err := r.GetUsersOrder(ctx, userID)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, op+": failed to get newly created order")
	}

	return order.ID, nil
}

func (r *Repository) checkOrderStatus(ctx context.Context, tx *sql.Tx, orderID uuid.UUID) (string, error) {
	const op = "repository.order.checkOrderStatus"

	var status string
	err := tx.QueryRowContext(ctx, "SELECT status FROM orders WHERE id = $1", orderID).Scan(&status)
	if err != nil {
		return "", errors.Wrap(err, op+": failed to get order status")
	}
	return status, nil
}

func (r *Repository) prepareStatements(ctx context.Context, tx *sql.Tx) (*sql.Stmt, *sql.Stmt, *sql.Stmt, *sql.Stmt, *sql.Stmt, *sql.Stmt, error) {
	stmtBookPrice, err := tx.PrepareContext(ctx, "SELECT price, stock FROM books WHERE id = $1 FOR UPDATE")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, errors.Wrap(err, ": failed to prepare book price statement")
	}

	stmtCheckExisting, err := tx.PrepareContext(ctx, "SELECT id, quantity FROM order_items WHERE order_id = $1 AND book_id = $2 FOR UPDATE")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, errors.Wrap(err, ": failed to prepare check existing item statement")
	}

	stmtUpdateItem, err := tx.PrepareContext(ctx, "UPDATE order_items SET quantity = quantity + $1 WHERE id = $2")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, errors.Wrap(err, ": failed to prepare update item statement")
	}

	stmtInsertItem, err := tx.PrepareContext(ctx, "INSERT INTO order_items (order_id, book_id, quantity, price) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, errors.Wrap(err, ": failed to prepare insert item statement")
	}

	stmtUpdateStock, err := tx.PrepareContext(ctx, "UPDATE books SET stock = stock - $1 WHERE id = $2")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, errors.Wrap(err, ": failed to prepare update stock statement")
	}

	stmtUpdateTotal, err := tx.PrepareContext(ctx, "UPDATE orders SET total_price = total_price + $1 WHERE id = $2")
	if err != nil {
		return nil, nil, nil, nil, nil, nil, errors.Wrap(err, ": failed to prepare update total statement")
	}

	return stmtBookPrice, stmtCheckExisting, stmtUpdateItem, stmtInsertItem, stmtUpdateStock, stmtUpdateTotal, nil
}

func (r *Repository) processItem(ctx context.Context, tx *sql.Tx, item orderModels.OrderItem, orderID uuid.UUID, stmtBookPrice, stmtCheckExisting, stmtUpdateItem, stmtInsertItem, stmtUpdateStock, stmtUpdateTotal *sql.Stmt) error {
	const op = "repository.order.processItem"

	var bookPrice float64
	var currentStock int
	err := stmtBookPrice.QueryRowContext(ctx, item.BookID).Scan(&bookPrice, &currentStock)
	if err != nil {
		return errors.Wrap(err, op+": failed to get book price and stock")
	}

	if currentStock < item.Quantity {
		return errors.Wrap(errors.New("insufficient stock"), op)
	}

	var existingItemID uuid.UUID
	var existingQuantity int
	err = stmtCheckExisting.QueryRowContext(ctx, orderID, item.BookID).Scan(&existingItemID, &existingQuantity)

	if err == sql.ErrNoRows {
		_, err = stmtInsertItem.ExecContext(ctx, orderID, item.BookID, item.Quantity, bookPrice)
		if err != nil {
			return errors.Wrap(err, op+": failed to add order item")
		}
	} else if err == nil {
		_, err = stmtUpdateItem.ExecContext(ctx, item.Quantity, existingItemID)
		if err != nil {
			return errors.Wrap(err, op+": failed to update order item quantity")
		}
	} else {
		return errors.Wrap(err, op+": failed to check existing item")
	}

	_, err = stmtUpdateStock.ExecContext(ctx, item.Quantity, item.BookID)
	if err != nil {
		return errors.Wrap(err, op+": failed to update book stock")
	}

	_, err = stmtUpdateTotal.ExecContext(ctx, bookPrice*float64(item.Quantity), orderID)
	if err != nil {
		return errors.Wrap(err, op+": failed to update order total price")
	}

	return nil
}

func (r *Repository) CreateUserOrderWithStatus(ctx context.Context, tx *sql.Tx, userID, status string) error {
	const op = "repository.order.CreateUserOrderWithStatus"

	_, err := tx.ExecContext(ctx, `
        INSERT INTO orders (user_id, status, total_price) 
        VALUES ($1, $2, 0)`, userID, status)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repository) InsertOrderItem(ctx context.Context, tx *sql.Tx, orderID, bookID uuid.UUID, quantity int, price float64) error {
	const op = "repository.order.InsertOrderItem"

	_, err := tx.ExecContext(ctx, `
        INSERT INTO order_items (order_id, book_id, quantity, price) 
        VALUES ($1, $2, $3, $4)`, orderID, bookID, quantity, price)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repository) GetOrderItemsFromOrderID(ctx context.Context, orderID string) (*[]orderModels.OrderItemFull, error) {
	const op = "repository.order.GetOrderItemsFromOrderID"

	orderUUID, err := uuid.FromString(orderID)
	if err != nil {
		return nil, errors.Wrap(err, op+": invalid order ID format")
	}

	stmt, err := r.DB.PrepareContext(ctx, `
        SELECT 
            oi.id AS order_item_id, 
            oi.book_id, 
            b.title AS book_title, 
            oi.quantity, 
            oi.price 
        FROM 
            order_items oi 
        JOIN 
            books b ON oi.book_id = b.id 
        JOIN 
            orders o ON oi.order_id = o.id 
        WHERE 
            oi.order_id = $1 
            AND o.status = 'draft'`)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, orderUUID)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer rows.Close()

	var orderItems []orderModels.OrderItemFull
	for rows.Next() {
		var orderItem orderModels.OrderItemFull
		err := rows.Scan(
			&orderItem.ID,
			&orderItem.BookID,
			&orderItem.Name,
			&orderItem.Quantity,
			&orderItem.Price,
		)
		if err != nil {
			return nil, errors.Wrap(err, op)
		}
		orderItems = append(orderItems, orderItem)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &orderItems, nil
}

func (r *Repository) RemoveCartItem(ctx context.Context, userID string, bookID uuid.UUID) error {
	const op = "repository.order.RemoveCartItem"
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, op+": failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var orderID uuid.UUID
	var quantity int
	var price float64
	err = tx.QueryRowContext(ctx, `
        SELECT o.id, oi.quantity, oi.price
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        WHERE o.user_id = $1 AND oi.book_id = $2 AND o.status = 'draft'
    `, userID, bookID).Scan(&orderID, &quantity, &price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(err, op+": no matching order or item found")
		}
		return errors.Wrap(err, op+": failed to get order details")
	}

	var bookExists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM books WHERE id = $1)", bookID).Scan(&bookExists)
	if err != nil {
		return errors.Wrap(err, op+": failed to check book existence")
	}
	if !bookExists {
		return errors.New("book not found")
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE books 
        SET stock = stock + $1 
        WHERE id = $2
    `, quantity, bookID)
	if err != nil {
		return errors.Wrap(err, op+": failed to update book stock")
	}

	res, err := tx.ExecContext(ctx, `
        DELETE FROM order_items 
        WHERE order_id = $1 AND book_id = $2
    `, orderID, bookID)
	if err != nil {
		return errors.Wrap(err, op+": failed to remove order item")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, op+": failed to check rows affected")
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected during deletion")
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE orders 
        SET total_price = total_price - $1 
        WHERE id = $2
    `, price*float64(quantity), orderID)
	if err != nil {
		return errors.Wrap(err, op+": failed to update order total")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, op+": failed to commit transaction")
	}

	return nil
}

func (r *Repository) AlterUserOrderByID(ctx context.Context, userID, orderID uuid.UUID) error {
	const op = "repository.order.AlterUserOrderByID"

	stmt, err := r.DB.PrepareContext(ctx, `
		UPDATE orders
		SET status = 'paid'
		WHERE user_id = $1
		  AND id = $2
		  AND status = 'draft'
	`)
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, userID, orderID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, op)
	}
	if rowsAffected == 0 {
		return errors.Wrap(errors.New("no draft order found for this user and order ID"), op)
	}

	return nil
}

func (r *Repository) ChangeStatusOfCart(ctx context.Context, userId, orderId string) error {
	const op = "repository.order.ChangeStatusOfCart"

	stmt, err := r.DB.PrepareContext(ctx, "UPDATE orders SET status = 'paid' WHERE user_id = $1 AND id = $2")
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userId, orderId)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repository) GetOrdersByUserID(ctx context.Context, userID string) ([]orderModels.HistoryOrderItem, error) {
	const op = "repository.order.GetOrdersByUserID"

	query := `
    WITH order_details AS (
        SELECT 
            o.id, 
            o.total_price,
            o.status,
            o.created_at,
            b.title as name, 
            oi.book_id, 
            oi.quantity, 
            oi.price
        FROM orders o
        JOIN order_items oi ON o.id = oi.order_id
        JOIN books b ON oi.book_id = b.id
        WHERE o.user_id = $1
    )
    SELECT 
        id, 
        total_price, 
        status, 
        created_at,
        json_agg(json_build_object(
            'book_id', book_id,
            'name', name,
            'quantity', quantity,
            'price', price
        )) as items
    FROM order_details
    GROUP BY id, total_price, status, created_at
    ORDER BY created_at DESC
    `

	var orders []orderModels.HistoryOrderItem
	var itemsJSON []byte

	rows, err := r.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.Wrap(err, op+": failed to execute query")
	}
	defer rows.Close()

	for rows.Next() {
		var order orderModels.HistoryOrderItem
		if err := rows.Scan(&order.ID, &order.TotalPrice, &order.Status, &order.CreatedAt, &itemsJSON); err != nil {
			return nil, errors.Wrap(err, op+": failed to scan row")
		}

		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, errors.Wrap(err, op+": failed to unmarshal items")
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, op+": rows iteration error")
	}

	return orders, nil
}
