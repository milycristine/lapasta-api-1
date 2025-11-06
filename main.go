package main

import (
	"embed"
	"flag"
	"log"

	"lapasta/config"
	dbsql "lapasta/database"
	utils "lapasta/internal/Utils"
	"lapasta/pkg/server"
)

var embeddedAssets embed.FS

func main() {
	var createConfig bool
	flag.BoolVar(&createConfig, "config", false, "cria o arquivo config.yaml")
	flag.Parse()

	if createConfig {
		config.CreateConfigFile()
		return
	}

	log.Println("Carregando config.yaml ...")
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	log.Println("Conectando ao SQL Server ...")
	connection, err := dbsql.MakeSQL(
		config.Yml.SQL.Host,
		config.Yml.SQL.Port,
		config.Yml.SQL.User,
		config.Yml.SQL.Password,
	)
	if err != nil {
		log.Fatal(err)
	}

	utils.SetSQLConn(connection)

	if err := server.StartServer(config.Yml.API.Port, connection, embeddedAssets); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
