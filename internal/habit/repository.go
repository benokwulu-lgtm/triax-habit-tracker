package habit

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll(ctx context.Context, userID int) ([]Habit, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, name, description, frequency, is_active
		FROM habits
		WHERE user_id = $1
		ORDER BY id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []Habit{}
	for rows.Next() {
		var h Habit
		var description sql.NullString
		if err := rows.Scan(&h.ID, &h.UserID, &h.Name, &description, &h.Frequency, &h.IsActive); err != nil {
			return nil, err
		}
		h.Description = description.String
		habits = append(habits, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return habits, nil
}

func (r *Repository) Create(ctx context.Context, h Habit) (Habit, error) {
	query := `
		INSERT INTO habits (user_id, name, description, frequency, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, h.UserID, h.Name, h.Description, h.Frequency, h.IsActive).Scan(&h.ID)
	if err != nil {
		return Habit{}, err
	}
	return h, nil
}

func (r *Repository) Get(ctx context.Context, id int, userID int) (Habit, error) {
	var h Habit
	var description sql.NullString
	query := `
		SELECT id, user_id, name, description, frequency, is_active
		FROM habits WHERE id = $1 AND user_id = $2
	`
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(&h.ID, &h.UserID, &h.Name, &description, &h.Frequency, &h.IsActive)
	if err == sql.ErrNoRows {
		return Habit{}, ErrNotFound
	}
	if err != nil {
		return Habit{}, err
	}
	h.Description = description.String
	return h, nil
}

func (r *Repository) Update(ctx context.Context, id int, userID int, h Habit) (Habit, error) {
	query := `
		UPDATE habits
		SET name = $1, description = $2, frequency = $3, is_active = $4
		WHERE id = $5 AND user_id = $6
	`
	result, err := r.db.ExecContext(ctx, query, h.Name, h.Description, h.Frequency, h.IsActive, id, userID)
	if err != nil {
		return Habit{}, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Habit{}, err
	}
	if rowsAffected == 0 {
		return Habit{}, ErrNotFound
	}
	h.ID = id
	return h, nil
}

func (r *Repository) Delete(ctx context.Context, id int, userID int) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM habits WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
