package infrastructure

import (
	"database/sql"
	"fmt"
	"hackathon-backend/settings_loader"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection(settings *settings_loader.SettingsLoader) (*sql.DB, error) {
	dbURL := settings.GetDatabaseURL()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com banco:  %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(settings.Database.MaxOpenConns)
	db.SetMaxIdleConns(settings.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(settings.Database.ConnMaxLifetime) * time.Second)

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao pingar banco de dados: %w", err)
	}

	log.Println("✅ Conexão com banco de dados estabelecida")
	return db, nil
}
