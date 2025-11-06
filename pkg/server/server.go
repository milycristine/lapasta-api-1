package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	dbsql "lapasta/database" 
)

/// StartServer inicializa o ServeMux, registra rotas e inicia o ListenAndServe.
/// - port: porta sem ":" (ex: "8080")
/// - db: instância de *dbsql.SQLStr
/// - assets: embed.FS (opcional) para arquivos estáticos
func StartServer(port string, db *dbsql.SQLStr, assets embed.FS) error {
	mux := http.NewServeMux()

	SetupRoutes(mux, db)

	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./public/images"))))

	if assets != (embed.FS{}) {
		mux.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.FS(assets))))
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Iniciando servidor na porta %s ", addr)
	return http.ListenAndServe(addr, mux)
}
