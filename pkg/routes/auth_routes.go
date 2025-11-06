package routes

import (
	"net/http"

	dbsql "lapasta/database"
	auth "lapasta/internal/AUTH"
)

/// RegisterAuthRoutes configura rotas relacionadas à autenticação.
/// - mux: ServeMux onde as rotas serão registradas.
/// - db: *dbsql.SQLStr injetado.
func RegisterAuthRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := auth.NewAuthRepository(db)
	svc := auth.NewAuthService(repo)

	mux.HandleFunc("/login", auth.LoginHandler(svc))
}
