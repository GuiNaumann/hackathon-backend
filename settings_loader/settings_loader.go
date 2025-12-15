package settings_loader

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type SettingsLoader struct {
	App      AppConfig
	Database DatabaseConfig
	Security SecurityConfig
	SMTP     SMTPConfig
	Storage  StorageConfig
}

type AppConfig struct {
	Port        string `toml:"port"`
	Environment string `toml:"environment"`
}

type DatabaseConfig struct {
	Host            string `toml:"host"`
	Port            string `toml:"port"`
	User            string `toml:"user"`
	Password        string `toml:"password"`
	DBName          string `toml:"dbname"`
	SSLMode         string `toml:"sslmode"`
	MaxOpenConns    int    `toml:"max_open_conns"`
	MaxIdleConns    int    `toml:"max_idle_conns"`
	ConnMaxLifetime int    `toml:"conn_max_lifetime"`
}

type SecurityConfig struct {
	CookieEncryptionKey string `toml:"cookie_encryption_key"`
	JWTSecret           string `toml:"jwt_secret"`
	CookieDomain        string `toml:"cookie_domain"`
	CookieSecure        bool   `toml:"cookie_secure"`
	CookieHTTPOnly      bool   `toml:"cookie_http_only"`
}

type SMTPConfig struct {
	Host      string `toml:"host"`
	Port      int    `toml:"port"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	FromEmail string `toml:"from_email"`
	FromName  string `toml:"from_name"`
}

type StorageConfig struct {
	RootPath    string `toml:"root_path"`
	UploadsPath string `toml:"uploads_path"`
	TempPath    string `toml:"temp_path"`
}

func NewSettingsLoader() *SettingsLoader {
	var settings SettingsLoader

	configPath := "./settings.toml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Arquivo settings.toml não encontrado em: %s", configPath)
	}

	if _, err := toml.DecodeFile(configPath, &settings); err != nil {
		log.Fatalf("Erro ao carregar settings.toml: %v", err)
	}

	// Validar campos obrigatórios
	if err := settings.Validate(); err != nil {
		log.Fatalf("Configuração inválida: %v", err)
	}

	log.Println("✅ Configurações carregadas com sucesso")
	return &settings
}

func (s *SettingsLoader) Validate() error {
	if s.App.Port == "" {
		return fmt.Errorf("app.port é obrigatório")
	}
	if s.Database.Host == "" {
		return fmt.Errorf("database. host é obrigatório")
	}
	if s.Database.DBName == "" {
		return fmt.Errorf("database.dbname é obrigatório")
	}
	if s.Security.CookieEncryptionKey == "" {
		return fmt.Errorf("security.cookie_encryption_key é obrigatório")
	}
	if s.Security.JWTSecret == "" {
		return fmt.Errorf("security.jwt_secret é obrigatório")
	}
	return nil
}

func (s *SettingsLoader) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		s.Database.User,
		s.Database.Password,
		s.Database.Host,
		s.Database.Port,
		s.Database.DBName,
		s.Database.SSLMode,
	)
}

func (s *SettingsLoader) GetAppPort() string {
	return s.App.Port
}

func (s *SettingsLoader) IsProduction() bool {
	return s.App.Environment == "production"
}
