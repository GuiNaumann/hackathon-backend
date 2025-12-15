package infrastructure

import (
	"database/sql"
	"fmt"
	"hackathon-backend/domain/usecases/usecase_impl"
	"hackathon-backend/infrastructure/modules/impl"
	repository_impl "hackathon-backend/infrastructure/repositories/impl"
	"hackathon-backend/settings_loader"
	"log"

	"github.com/gorilla/mux"
)

type SetupConfig struct {
	DB                *sql.DB
	Settings          *settings_loader.SettingsLoader
	AuthRepository    *repository_impl.AuthRepositoryImpl
	PermRepository    *repository_impl.PermissionRepositoryImpl
	AuthUseCase       *usecase_impl.AuthUseCaseImpl
	PermissionUseCase *usecase_impl.PermissionUseCaseImpl
	UserCrudUseCase   *usecase_impl.UserCrudUseCaseImpl
}

func Setup(router *mux.Router, settings *settings_loader.SettingsLoader) (*SetupConfig, error) {
	log.Println("🔧 Iniciando setup da aplicação...")

	db, err := NewDatabaseConnection(settings)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	authRepository := repository_impl.NewAuthRepositoryImpl(db)
	permRepository := repository_impl.NewPermissionRepositoryImpl(db)

	authUseCase := usecase_impl.NewAuthUseCaseImpl(authRepository, settings)
	permUseCase := usecase_impl.NewPermissionUseCaseImpl(permRepository, authRepository)
	userCrudUseCase := usecase_impl.NewUserCrudUseCaseImpl(authRepository, permRepository)

	authModule := module_impl.NewAuthModule(authUseCase, settings)
	permModule := module_impl.NewPermissionModule(permUseCase)
	userCrudModule := module_impl.NewUserCrudModule(userCrudUseCase)
	healthModule := module_impl.NewHealthModule()

	// 🔹 ROOT /api
	apiRouter := router.PathPrefix("/api").Subrouter()

	// 🔓 Públicas
	authModule.RegisterPublicRoutes(apiRouter)
	healthModule.RegisterRoutes(apiRouter)

	// 🔐 Privadas → /api/private/*
	privateRouter := apiRouter.PathPrefix("/private").Subrouter()

	privateRouter.Use(NewAuthMiddleware(authRepository, settings))
	privateRouter.Use(NewPermissionMiddleware(permUseCase))

	authModule.RegisterPrivateRoutes(privateRouter)
	permModule.RegisterRoutes(privateRouter)
	userCrudModule.RegisterRoutes(privateRouter)

	log.Println("✅ Setup concluído com sucesso")

	return &SetupConfig{
		DB:                db,
		Settings:          settings,
		AuthRepository:    authRepository,
		PermRepository:    permRepository,
		AuthUseCase:       authUseCase,
		PermissionUseCase: permUseCase,
		UserCrudUseCase:   userCrudUseCase,
	}, nil
}

func (s *SetupConfig) CloseDB() {
	if s.DB != nil {
		s.DB.Close()
		log.Println("✅ Conexão com banco de dados fechada")
	}
}
