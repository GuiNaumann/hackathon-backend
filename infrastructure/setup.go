package infrastructure

import (
	"database/sql"
	"fmt"
	"hackathon-backend/domain/usecases/usecase_impl"
	"hackathon-backend/infrastructure/modules/impl"
	"hackathon-backend/infrastructure/repositories/impl"
	"hackathon-backend/settings_loader"
	"log"

	"github.com/gorilla/mux"
)

type SetupConfig struct {
	DB             *sql.DB
	Settings       *settings_loader.SettingsLoader
	AuthRepository *repository_impl.AuthRepositoryImpl
	AuthUseCase    *usecase_impl.AuthUseCaseImpl
}

func Setup(router *mux.Router, settings *settings_loader.SettingsLoader) (*SetupConfig, error) {
	log.Println("🔧 Iniciando setup da aplicação...")

	// 1. Conectar ao banco de dados
	db, err := NewDatabaseConnection(settings)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	// 2. Inicializar Repositories
	authRepository := repository_impl.NewAuthRepositoryImpl(db)

	// 3. Inicializar UseCases
	authUseCase := usecase_impl.NewAuthUseCaseImpl(authRepository, settings)

	// 4. Inicializar Módulos HTTP
	authModule := module_impl.NewAuthModule(authUseCase, settings)
	healthModule := module_impl.NewHealthModule()

	// 5. Registrar Rotas Públicas (sem autenticação)
	publicRouter := router.PathPrefix("/api").Subrouter()
	authModule.RegisterPublicRoutes(publicRouter)
	healthModule.RegisterRoutes(publicRouter)

	// 6. Registrar Rotas Privadas (com autenticação)
	privateRouter := router.PathPrefix("/private").Subrouter()
	privateRouter.Use(NewAuthMiddleware(authRepository, settings))
	authModule.RegisterPrivateRoutes(privateRouter)

	log.Println("✅ Setup concluído com sucesso")

	return &SetupConfig{
		DB:             db,
		Settings:       settings,
		AuthRepository: authRepository,
		AuthUseCase:    authUseCase,
	}, nil
}

func (s *SetupConfig) CloseDB() {
	if s.DB != nil {
		s.DB.Close()
		log.Println("✅ Conexão com banco de dados fechada")
	}
}
