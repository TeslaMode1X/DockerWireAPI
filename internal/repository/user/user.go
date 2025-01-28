package user

import (
	"context"
	"database/sql"
	"fmt"
	userModel "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/user"
	"github.com/TeslaMode1X/DockerWireAPI/internal/repository"
	"github.com/gofrs/uuid"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) CheckUserExists(ctx context.Context, username string) (bool, error) {
	const op = "repo.user.CheckUserExists"

	stmt, err := r.DB.PrepareContext(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1);")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var exists bool

	err = stmt.QueryRowContext(ctx, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id uuid.UUID) (userModel.User, error) {
	const op = "repo.user.FindUserByID"

	stmt, err := r.DB.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
	if err != nil {
		return userModel.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
	}
	defer stmt.Close()

	var user userModel.User
	err = stmt.QueryRowContext(ctx, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return userModel.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
