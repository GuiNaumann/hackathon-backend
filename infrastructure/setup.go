package infrastructure

import (
	"database/sql"
	"fmt"
	"hackathon-backend/infrastructure/repositories"
	"log"

	"hackathon-backend/domain/usecases/usecase_impl"
	"hackathon-backend/infrastructure/modules/impl"
	repository_impl "hackathon-backend/infrastructure/repositories/impl"
	"hackathon-backend/settings_loader"

	"github.com/gorilla/mux"
)

type SetupConfig struct {
	DB                          *sql.DB
	Settings                    *settings_loader.SettingsLoader
	AuthRepository              *repository_impl.AuthRepositoryImpl
	PermRepository              *repositories.PermissionRepositoryImpl
	InitiativeRepository        *repository_impl.InitiativeRepositoryImpl
	CommentRepository           *repository_impl.CommentRepositoryImpl
	AuthUseCase                 *usecase_impl.AuthUseCaseImpl
	PermissionUseCase           *usecase_impl.PermissionUseCaseImpl
	UserCrudUseCase             *usecase_impl.UserCrudUseCaseImpl
	InitiativeUseCase           *usecase_impl.InitiativeUseCaseImpl
	CommentUseCase              *usecase_impl.CommentUseCaseImpl
	InitiativeHistoryRepository *repository_impl.InitiativeHistoryRepositoryImpl
	InitiativeHistoryUseCase    *usecase_impl.InitiativeHistoryUseCaseImpl
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
	permRepository := repositories.NewPermissionRepositoryImpl(db)
	initiativeRepository := repository_impl.NewInitiativeRepositoryImpl(db)
	commentRepository := repository_impl.NewCommentRepositoryImpl(db)
	initiativeHistoryRepository := repository_impl.NewInitiativeHistoryRepositoryImpl(db)
	cancellationRepository := repository_impl.NewCancellationRepositoryImpl(db)

	// 3. Inicializar UseCases
	authUseCase := usecase_impl.NewAuthUseCaseImpl(authRepository, settings)
	permUseCase := usecase_impl.NewPermissionUseCaseImpl(permRepository, authRepository)
	userCrudUseCase := usecase_impl.NewUserCrudUseCaseImpl(authRepository, permRepository)
	initiativeUseCase := usecase_impl.NewInitiativeUseCaseImpl(initiativeRepository, permRepository)
	commentUseCase := usecase_impl.NewCommentUseCaseImpl(commentRepository, initiativeRepository, permRepository)
	initiativeHistoryUseCase := usecase_impl.NewInitiativeHistoryUseCaseImpl(initiativeHistoryRepository)
	cancellationUseCase := usecase_impl.NewCancellationUseCaseImpl(cancellationRepository, initiativeRepository, permRepository)

	// 4. Inicializar Módulos HTTP
	authModule := module_impl.NewAuthModule(authUseCase, settings)
	permModule := module_impl.NewPermissionModule(permUseCase)
	userCrudModule := module_impl.NewUserCrudModule(userCrudUseCase)
	initiativeModule := module_impl.NewInitiativeModule(initiativeUseCase, initiativeHistoryUseCase, cancellationUseCase)
	commentModule := module_impl.NewCommentModule(commentUseCase)
	healthModule := module_impl.NewHealthModule()

	// 5. Registrar Rotas Públicas (sem autenticação)
	publicRouter := router.PathPrefix("/api").Subrouter()
	authModule.RegisterPublicRoutes(publicRouter)
	healthModule.RegisterRoutes(publicRouter)

	// 6. Registrar Rotas Privadas (com autenticação + permissões)
	privateRouter := router.PathPrefix("/api/private").Subrouter()

	// Middleware:  Autenticação
	privateRouter.Use(NewAuthMiddleware(authRepository, settings))

	// Middleware: Verificação de permissões
	privateRouter.Use(NewPermissionMiddleware(permUseCase))

	// Registrar rotas privadas
	authModule.RegisterPrivateRoutes(privateRouter)
	permModule.RegisterRoutes(privateRouter)
	userCrudModule.RegisterRoutes(privateRouter)
	initiativeModule.RegisterRoutes(privateRouter)
	commentModule.RegisterRoutes(privateRouter)

	log.Println("✅ Setup concluído com sucesso")

	return &SetupConfig{
		DB:                   db,
		Settings:             settings,
		AuthRepository:       authRepository,
		PermRepository:       permRepository,
		InitiativeRepository: initiativeRepository,
		CommentRepository:    commentRepository,
		AuthUseCase:          authUseCase,
		PermissionUseCase:    permUseCase,
		UserCrudUseCase:      userCrudUseCase,
		InitiativeUseCase:    initiativeUseCase,
		CommentUseCase:       commentUseCase,
	}, nil
}

func (s *SetupConfig) CloseDB() {
	if s.DB != nil {
		s.DB.Close()
		log.Println("✅ Conexão com banco de dados fechada")
	}
}
