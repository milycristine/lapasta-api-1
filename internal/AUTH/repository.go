package auth

import (
	"database/sql"
	"fmt"

	database "lapasta/database"
	"lapasta/internal/models"
)

type AuthRepository interface {
	Autenticar(username string) (models.Login, error)
	GetUsuario(username string) (models.User, error)
}
type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(conn *database.SQLStr) AuthRepository {
	return &authRepository{
		db: conn.DB(),
	}
}
func (r *authRepository) Autenticar(username string) (models.Login, error) {
	var login models.Login

	query := `
		SELECT username, password 
		FROM AUTH WITH (NOLOCK) 
		WHERE username = @Username
	`

	err := r.db.QueryRow(query, sql.Named("Username", username)).
		Scan(&login.Username, &login.PasswordCriptografado)

	if err != nil {
		if err == sql.ErrNoRows {
			return login, fmt.Errorf("usuário não encontrado")
		}
		return login, err
	}

	return login, nil
}
func (r *authRepository) GetUsuario(email string) (models.User, error) {
	var user models.User

	queryFunc := `
		SELECT ID, Email, Admin, Nome, Sobrenome 
		FROM Funcionarios WITH (NOLOCK) 
		WHERE Email = @Email
	`

	var admin int
	err := r.db.QueryRow(queryFunc, sql.Named("Email", email)).
		Scan(&user.ID, &user.Email, &admin, &user.Nome, &user.Sobrenome)

	if err == nil {
		user.Admin = admin
		if admin == 1 {
			user.Tipo = "admin"
		} else {
			user.Tipo = "funcionario"
		}
		return user, nil
	}

	queryMot := `
		SELECT Id, Email, Nome 
		FROM Motoristas WITH (NOLOCK) 
		WHERE Email = @Email
	`

	err = r.db.QueryRow(queryMot, sql.Named("Email", email)).
		Scan(&user.ID, &user.Email, &user.Nome)

	if err == nil {
		user.Tipo = "motorista"
		return user, nil
	}

	return user, fmt.Errorf("usuário não encontrado")
}
