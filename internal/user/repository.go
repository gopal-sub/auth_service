package user

import (
	"database/sql"
	"errors"
)


type Repository struct{
	db *sql.DB
}


type UserRepository interface{
	FindUserByEmail(email string) (User, error)
	Create(User) error
	VerifyAccount(email string) (bool, error)
	ResetPassword(email string, passwordHash string) error
}

func NewRepository(db *sql.DB) *Repository{
	return &Repository{
		db: db,
	}
}


func (r *Repository) Create(user User) error{


	query := `
		INSERT INTO users (id, email, password_hash, created_at, is_verified)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.IsVerified,
	)
	return err


}

func (r *Repository) FindUserByEmail(email string)(User, error) {
	query := `
	SELECT
		id,
		email,
		password_hash,
		created_at,
		is_verified
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
		&user.IsVerified,

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


func (r *Repository) VerifyAccount(email string) (accountVerified bool, err error){
	query := `
	UPDATE users
	SET	is_verified = TRUE
	WHERE email = $1;
	`
	result, err := r.db.Exec(query, email)
	if err != nil{
		return false, err
	}
	rowsAfffected, err := result.RowsAffected()
	if rowsAfffected == 0 {
		return false, UserDoesNotExistErr
	}
	return true, nil

}


func (r *Repository) ResetPassword(email string, passwordHash string) error{
	query := `
	UPDATE users
	SET	password_hash = $2
	WHERE email = $1;
	`
	_, err := r.db.Exec(query, email, passwordHash)
	return err
}


