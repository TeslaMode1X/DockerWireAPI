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

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, user_id, status, total_price WHERE user_id = ? FROM order")
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

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, user_id, status, total_price WHERE order_id = ? FROM order")
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

	stmt, err := r.DB.PrepareContext(ctx, "INSERT INTO order(user_id) VALUES (?)")
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

	stmt, err := r.DB.PrepareContext(ctx, "SELECT EXISTS(SELECT 1 FROM orders WHERE user_id = ? AND status = 'draft')")
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

	stmtBookPrice, err := tx.PrepareContext(ctx, "SELECT price FROM books WHERE id = ?")
	if err != nil {
		return errors.Wrap(err, op+": failed to prepare book price statement")
	}
	defer stmtBookPrice.Close()

	stmtInsertItem, err := tx.PrepareContext(ctx, "INSERT INTO order_items (order_id, book_id, quantity, price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, op+": failed to prepare insert order item statement")
	}
	defer stmtInsertItem.Close()

	stmtUpdateTotal, err := tx.PrepareContext(ctx, "UPDATE orders SET total_price = total_price + ? WHERE id = ?")
	if err != nil {
		return errors.Wrap(err, op+": failed to prepare update order total statement")
	}
	defer stmtUpdateTotal.Close()

	for _, item := range *items {
		var bookPrice float64
		err = stmtBookPrice.QueryRowContext(ctx, item.BookID).Scan(&bookPrice)
		if err != nil {
			return errors.Wrap(err, op+": failed to get book price")
		}

		_, err = stmtInsertItem.ExecContext(ctx, orderID, item.BookID, item.Quantity, bookPrice)
		if err != nil {
			return errors.Wrap(err, op+": failed to add order item")
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
