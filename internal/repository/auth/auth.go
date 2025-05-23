package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	model "github.com/TeslaMode1X/DockerWireAPI/internal/domain/models/auth"
	"github.com/TeslaMode1X/DockerWireAPI/internal/repository"
	"github.com/gofrs/uuid"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) Register(ctx context.Context, user model.Registration) (uuid.UUID, error) {
	const op = "repo.user.Register"

	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := r.DB.PrepareContext(ctx, "INSERT INTO users (id, username, email, password, created_at) VALUES($1, $2, $3, $4, $5)")
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id, user.Username, user.Email, user.Password, time.Now())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repository) Login(ctx context.Context, user model.Login) (uuid.UUID, int, error) {
	const op = "repo.user.Login"

	var storedPassword string
	var userID uuid.UUID
	var role int

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, password, role FROM users where email = $1")
	if err != nil {
		return uuid.Nil, 0, err
	}

	err = stmt.QueryRowContext(ctx, user.Email).Scan(&userID, &storedPassword, &role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, 0, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return uuid.Nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	hashedPassword := sha256.Sum256([]byte(user.Password))
	hashedPasswordHex := hex.EncodeToString(hashedPassword[:])

	if storedPassword != hashedPasswordHex {
		return uuid.Nil, 0, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
	}

	return userID, role, nil
}
