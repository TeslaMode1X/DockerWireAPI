package order

import (
	"context"
	"database/sql"
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

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, user_id, status, total_price FROM orders WHERE user_id = $1")
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

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, user_id, status, total_price WHERE order_id = $1 FROM order")
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

	exists, err := r.CheckOrderExists(ctx, userID)
	if err != nil {
		return errors.Wrap(err, op+": failed to check existing order")
	}

	var orderID uuid.UUID
	if exists {
		order, err := r.GetUsersOrder(ctx, userID)
		if err != nil {
			return errors.Wrap(err, op+": failed to get existing order")
		}
		orderID = order.ID
	} else {
		err := r.CreateUserOrder(ctx, userID)
		if err != nil {
			return errors.Wrap(err, op+": failed to create new order")
		}

		order, err := r.GetUsersOrder(ctx, userID)
		if err != nil {
			return errors.Wrap(err, op+": failed to get newly created order")
		}
		orderID = order.ID
	}

	stmtBookPrice, err := tx.PrepareContext(ctx, "SELECT price, stock FROM books WHERE id = $1 FOR UPDATE")
	if err != nil {
		return errors.Wrap(err, op+": failed to prepare book price statement")
	}
	defer stmtBookPrice.Close()

	stmtCheckExisting, err := tx.PrepareContext(ctx, `
        SELECT id, quantity FROM order_items 
        WHERE order_id = $1 AND book_id = $2
        FOR UPDATE`)
	if err != nil {
		return errors.Wrap(err, op+": failed to prepare check existing item statement")
	}
	defer stmtCheckExisting.Close()

	stmtUpdateItem, err := tx.PrepareContext(ctx, `
        UPDATE order_items 
        SET quantity = quantity + $1 
        WHERE id = $2`)
	if err != nil {
		return errors.Wrap(err, op+": failed to prepare update item statement")
	}
	defer stmtUpdateItem.Close()

	stmtInsertItem, err := tx.PrepareContext(ctx, `
        INSERT INTO order_items (order_id, book_id, quantity, price) 
        VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return errors.Wrap(err, op+": failed to prepare insert item statement")
	}
	defer stmtInsertItem.Close()

	stmtUpdateStock, err := tx.PrepareContext(ctx, `
        UPDATE books 
        SET stock = stock - $1 
        WHERE id = $2`)
	if err != nil {
		return errors.Wrap(err, op+": failed to prepare update stock statement")
	}
	defer stmtUpdateStock.Close()

	stmtUpdateTotal, err := tx.PrepareContext(ctx, `
        UPDATE orders 
        SET total_price = total_price + $1 
        WHERE id = $2`)
	if err != nil {
		return errors.Wrap(err, op+": failed to prepare update total statement")
	}
	defer stmtUpdateTotal.Close()

	for _, item := range *items {
		var bookPrice float64
		var currentStock int
		err = stmtBookPrice.QueryRowContext(ctx, item.BookID).Scan(&bookPrice, &currentStock)
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
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, op+": failed to commit transaction")
	}

	return nil
}

func (r *Repository) GetOrderItemsFromOrderID(ctx context.Context, orderID string) (*[]orderModels.OrderItemFull, error) {
	const op = "repository.order.GetOrderItemsFromOrderID"

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, book_id, quantity, price FROM order_items WHERE order_id = $1")
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, orderID)
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
		return errors.Wrap(err, op+": failed to get order details")
	}

	_, err = tx.ExecContext(ctx, `
        UPDATE books 
        SET stock = stock + $1 
        WHERE id = $2
    `, quantity, bookID)
	if err != nil {
		return errors.Wrap(err, op+": failed to update book stock")
	}

	_, err = tx.ExecContext(ctx, `
        DELETE FROM order_items 
        WHERE order_id = $1 AND book_id = $2
    `, orderID, bookID)
	if err != nil {
		return errors.Wrap(err, op+": failed to remove order item")
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
