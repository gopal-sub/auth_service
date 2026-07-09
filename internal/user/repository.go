package user

import (
	"database/sql"
	"errors"
)


type Repository struct{
	db *sql.DB
}
func NewRepository(db *sql.DB) *Repository{
	return &Repository{
		db: db,
	}
}


func (r *Repository) Create(user User) error{


	query := `
		INSERT INTO users (id, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	)
	return err


}

func (r *Repository) FindUserByEmail(email string)(User, error) {
	query := `
	SELECT
		id,
		email,
		password_hash,
		created_at
	FROM users
	WHERE email = $1;
	`
	var user User

	row := r.db.QueryRow(query, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,

	)
	if errors.Is(err, sql.ErrNoRows){
		// there is no row for in the database for the search
		return User{}, sql.ErrNoRows
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
	
}




