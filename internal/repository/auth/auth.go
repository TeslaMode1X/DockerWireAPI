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

	stmt, err := r.DB.PrepareContext(ctx, "INSERT INTO users (username, password, role, created_at) VALUES($1, $2, $3, $4)")
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.Username, user.Password, user.Role, time.Now())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repository) Login(ctx context.Context, user model.Login) (uuid.UUID, error) {
	const op = "repo.user.Login"

	var storedPassword string
	var userID uuid.UUID

	stmt, err := r.DB.PrepareContext(ctx, "SELECT id, password FROM users where username = $1")
	if err != nil {
		return uuid.Nil, err
	}

	err = stmt.QueryRowContext(ctx, user.Username).Scan(&userID, &storedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	hashedPassword := sha256.Sum256([]byte(user.Password))
	hashedPasswordHex := hex.EncodeToString(hashedPassword[:])

	if storedPassword != hashedPasswordHex {
		return uuid.Nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
	}

	return userID, nil
}
