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
	InitiativeHistoryRepository *repository_impl.InitiativeHistoryRepositoryImpl
	CancellationRepository      *repository_impl.CancellationRepositoryImpl
	AIRepository                *repository_impl.AIRepositoryImpl
	SectorRepository            *repository_impl.SectorRepositoryImpl
	PrioritizationRepository    *repository_impl.PrioritizationRepositoryImpl // NOVO
	AuthUseCase                 *usecase_impl.AuthUseCaseImpl
	PermissionUseCase           *usecase_impl.PermissionUseCaseImpl
	UserCrudUseCase             *usecase_impl.UserCrudUseCaseImpl
	InitiativeUseCase           *usecase_impl.InitiativeUseCaseImpl
	CommentUseCase              *usecase_impl.CommentUseCaseImpl
	InitiativeHistoryUseCase    *usecase_impl.InitiativeHistoryUseCaseImpl
	CancellationUseCase         *usecase_impl.CancellationUseCaseImpl
	AIUseCase                   *usecase_impl.AIUseCaseImpl
	SectorUseCase               *usecase_impl.SectorUseCaseImpl
	PrioritizationUseCase       *usecase_impl.PrioritizationUseCaseImpl // NOVO
}

func Setup(router *mux.Router, settings *settings_loader.SettingsLoader) (*SetupConfig, error) {
	log.Println("🔧 Iniciando setup da aplicação...")

	// 1. Conectar ao banco de dados
	db, err := NewDatabaseConnection(settings)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco:  %w", err)
	}

	// 2. Inicializar Repositories
	log.Println("📦 Inicializando repositories...")
	authRepository := repository_impl.NewAuthRepositoryImpl(db)
	permRepository := repositories.NewPermissionRepositoryImpl(db)
	initiativeRepository := repository_impl.NewInitiativeRepositoryImpl(db)
	commentRepository := repository_impl.NewCommentRepositoryImpl(db)
	initiativeHistoryRepository := repository_impl.NewInitiativeHistoryRepositoryImpl(db)
	cancellationRepository := repository_impl.NewCancellationRepositoryImpl(db)
	aiRepository := repository_impl.NewAIRepositoryImpl(settings)
	sectorRepository := repository_impl.NewSectorRepositoryImpl(db)
	prioritizationRepository := repository_impl.NewPrioritizationRepositoryImpl(db) // NOVO

	// 3. Inicializar UseCases
	log.Println("⚙️  Inicializando use cases...")
	authUseCase := usecase_impl.NewAuthUseCaseImpl(authRepository, settings)
	permUseCase := usecase_impl.NewPermissionUseCaseImpl(permRepository, authRepository)
	userCrudUseCase := usecase_impl.NewUserCrudUseCaseImpl(authRepository, permRepository)
	initiativeUseCase := usecase_impl.NewInitiativeUseCaseImpl(
		initiativeRepository,
		initiativeHistoryRepository,
		permRepository,
		authRepository,
	)
	commentUseCase := usecase_impl.NewCommentUseCaseImpl(commentRepository, initiativeRepository, permRepository)
	initiativeHistoryUseCase := usecase_impl.NewInitiativeHistoryUseCaseImpl(initiativeHistoryRepository)

	cancellationUseCase := usecase_impl.NewCancellationUseCaseImpl(
		cancellationRepository,
		initiativeRepository,
		initiativeHistoryRepository,
		permRepository,
	)

	aiUseCase := usecase_impl.NewAIUseCaseImpl(aiRepository)
	sectorUseCase := usecase_impl.NewSectorUseCaseImpl(sectorRepository)

	// NOVO: PrioritizationUseCase
	prioritizationUseCase := usecase_impl.NewPrioritizationUseCaseImpl(
		prioritizationRepository,
		initiativeRepository,
		authRepository,
		permRepository,
		sectorRepository, // ADICIONAR
	)

	// 4. Inicializar Módulos HTTP
	log.Println("🌐 Inicializando módulos HTTP...")
	authModule := module_impl.NewAuthModule(authUseCase, settings)
	permModule := module_impl.NewPermissionModule(permUseCase)
	userCrudModule := module_impl.NewUserCrudModule(userCrudUseCase)
	initiativeModule := module_impl.NewInitiativeModule(initiativeUseCase, initiativeHistoryUseCase, cancellationUseCase)
	commentModule := module_impl.NewCommentModule(commentUseCase)
	aiModule := module_impl.NewAIModule(aiUseCase)
	sectorModule := module_impl.NewSectorModule(sectorUseCase)
	prioritizationModule := module_impl.NewPrioritizationModule(prioritizationUseCase) // NOVO
	healthModule := module_impl.NewHealthModule()

	// 5. Registrar Rotas Públicas
	log.Println("🔓 Registrando rotas públicas...")
	publicRouter := router.PathPrefix("/api").Subrouter()
	authModule.RegisterPublicRoutes(publicRouter)
	healthModule.RegisterRoutes(publicRouter)

	// 6. Registrar Rotas Privadas
	log.Println("🔒 Registrando rotas privadas...")
	privateRouter := router.PathPrefix("/api/private").Subrouter()
	privateRouter.Use(NewAuthMiddleware(authRepository, settings))
	privateRouter.Use(NewPermissionMiddleware(permUseCase))

	authModule.RegisterPrivateRoutes(privateRouter)
	permModule.RegisterRoutes(privateRouter)
	userCrudModule.RegisterRoutes(privateRouter)
	initiativeModule.RegisterRoutes(privateRouter)
	commentModule.RegisterRoutes(privateRouter)
	aiModule.RegisterRoutes(privateRouter)
	sectorModule.RegisterRoutes(privateRouter)
	prioritizationModule.RegisterRoutes(privateRouter) // NOVO

	log.Println("✅ Setup concluído com sucesso!")

	return &SetupConfig{
		DB:                          db,
		Settings:                    settings,
		AuthRepository:              authRepository,
		PermRepository:              permRepository,
		InitiativeRepository:        initiativeRepository,
		CommentRepository:           commentRepository,
		InitiativeHistoryRepository: initiativeHistoryRepository,
		CancellationRepository:      cancellationRepository,
		AIRepository:                aiRepository,
		SectorRepository:            sectorRepository,
		PrioritizationRepository:    prioritizationRepository,
		AuthUseCase:                 authUseCase,
		PermissionUseCase:           permUseCase,
		UserCrudUseCase:             userCrudUseCase,
		InitiativeUseCase:           initiativeUseCase,
		CommentUseCase:              commentUseCase,
		InitiativeHistoryUseCase:    initiativeHistoryUseCase,
		CancellationUseCase:         cancellationUseCase,
		AIUseCase:                   aiUseCase,
		SectorUseCase:               sectorUseCase,
		PrioritizationUseCase:       prioritizationUseCase,
	}, nil
}

func (s *SetupConfig) CloseDB() {
	if s.DB != nil {
		s.DB.Close()
		log.Println("✅ Conexão com banco de dados fechada")
	}
}
