package infrastructure

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"hackathon-backend/domain/usecases/usecase_impl"
	"hackathon-backend/infrastructure/modules/impl"
	repository_impl "hackathon-backend/infrastructure/repositories/impl"
	"hackathon-backend/settings_loader"
	"log"
)

type SetupConfig struct {
	DB                *sql.DB
	Settings          *settings_loader.SettingsLoader
	AuthRepository    *repository_impl.AuthRepositoryImpl
	PermRepository    *repository_impl.PermissionRepositoryImpl
	AuthUseCase       *usecase_impl.AuthUseCaseImpl
	PermissionUseCase *usecase_impl.PermissionUseCaseImpl
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
	permRepository := repository_impl.NewPermissionRepositoryImpl(db)

	// 3. Inicializar UseCases
	authUseCase := usecase_impl.NewAuthUseCaseImpl(authRepository, settings)
	permUseCase := usecase_impl.NewPermissionUseCaseImpl(permRepository, authRepository)

	// 4. Inicializar Módulos HTTP
	authModule := module_impl.NewAuthModule(authUseCase, settings)
	permModule := module_impl.NewPermissionModule(permUseCase)
	healthModule := module_impl.NewHealthModule()

	// 5. Registrar Rotas Públicas (sem autenticação)
	publicRouter := router.PathPrefix("/api").Subrouter()
	authModule.RegisterPublicRoutes(publicRouter)
	healthModule.RegisterRoutes(publicRouter)

	// 6. Registrar Rotas Privadas (com autenticação + permissões)
	privateRouter := router.PathPrefix("/private").Subrouter()

	// Primeiro middleware:  Autenticação
	privateRouter.Use(NewAuthMiddleware(authRepository, settings))

	// Segundo middleware: Verificação de permissões
	privateRouter.Use(NewPermissionMiddleware(permUseCase))

	// Registrar rotas privadas
	authModule.RegisterPrivateRoutes(privateRouter)
	permModule.RegisterRoutes(privateRouter)

	log.Println("✅ Setup concluído com sucesso")

	return &SetupConfig{
		DB:                db,
		Settings:          settings,
		AuthRepository:    authRepository,
		PermRepository:    permRepository,
		AuthUseCase:       authUseCase,
		PermissionUseCase: permUseCase,
	}, nil
}

func (s *SetupConfig) CloseDB() {
	if s.DB != nil {
		s.DB.Close()
		log.Println("✅ Conexão com banco de dados fechada")
	}
}
