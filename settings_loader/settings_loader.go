package settings_loader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type SettingsLoader struct {
	App      AppConfig
	Database DatabaseConfig
	Security SecurityConfig
	SMTP     SMTPConfig
	Storage  StorageConfig
	AI       AIConfig // NOVO
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

// NOVO: Configuração de IA
type AIConfig struct {
	GeminiAPIKey   string  `toml:"gemini_api_key"`
	GeminiModel    string  `toml:"gemini_model"`
	MaxTokens      int     `toml:"max_tokens"`
	Temperature    float64 `toml:"temperature"`
	RequestTimeout int     `toml:"request_timeout"` // em segundos
}

func NewSettingsLoader() *SettingsLoader {
	var settings SettingsLoader

	configPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configPath = filepath.Join(configPath, "settings.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Arquivo settings.toml não encontrado em: %s", configPath)
	}

	if _, err := toml.DecodeFile(configPath, &settings); err != nil {
		log.Fatalf("Erro ao carregar settings. toml: %v", err)
	}

	// Aplicar valores padrão para AI se não configurado
	settings.applyDefaults()

	// Validar campos obrigatórios
	if err := settings.Validate(); err != nil {
		log.Fatalf("Configuração inválida: %v", err)
	}

	log.Println("✅ Configurações carregadas com sucesso")
	return &settings
}

// NOVO: Aplicar valores padrão
func (s *SettingsLoader) applyDefaults() {
	// Defaults para AI
	if s.AI.GeminiModel == "" {
		s.AI.GeminiModel = "gemini-2.0-flash"
	}
	if s.AI.MaxTokens == 0 {
		s.AI.MaxTokens = 2000
	}
	if s.AI.Temperature == 0 {
		s.AI.Temperature = 0.7
	}
	if s.AI.RequestTimeout == 0 {
		s.AI.RequestTimeout = 30
	}
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
		return fmt.Errorf("security. cookie_encryption_key é obrigatório")
	}
	if s.Security.JWTSecret == "" {
		return fmt.Errorf("security.jwt_secret é obrigatório")
	}

	// NOVO: Validar AI (apenas aviso, não obrigatório)
	if s.AI.GeminiAPIKey == "" {
		log.Println("⚠️  Aviso: ai.gemini_api_key não configurado - funcionalidades de IA estarão desabilitadas")
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

// NOVO: Helpers para AI
func (s *SettingsLoader) GetGeminiAPIKey() string {
	return s.AI.GeminiAPIKey
}

func (s *SettingsLoader) GetGeminiModel() string {
	return s.AI.GeminiModel
}

func (s *SettingsLoader) GetAIMaxTokens() int {
	return s.AI.MaxTokens
}

func (s *SettingsLoader) GetAITemperature() float64 {
	return s.AI.Temperature
}

func (s *SettingsLoader) GetAIRequestTimeout() int {
	return s.AI.RequestTimeout
}

func (s *SettingsLoader) IsAIEnabled() bool {
	return s.AI.GeminiAPIKey != ""
}
