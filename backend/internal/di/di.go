package di

import (
	"context"
	"log"

	"github.com/joineroff/social-network/backend/internal/config"
	"github.com/joineroff/social-network/backend/internal/framework/routing"
	"github.com/joineroff/social-network/backend/internal/infrastructure/database"
	"github.com/joineroff/social-network/backend/internal/repository"
	"github.com/joineroff/social-network/backend/internal/repository/friendRepository"
	"github.com/joineroff/social-network/backend/internal/repository/profileRepository"
	"github.com/joineroff/social-network/backend/internal/repository/userRepository"
	"github.com/joineroff/social-network/backend/internal/service"
	"github.com/joineroff/social-network/backend/internal/usecase"
	"github.com/joineroff/social-network/backend/pkg/closer"
	"github.com/joineroff/social-network/backend/pkg/server"
)

type DIContainer struct {
	Config *config.Config
	closer *closer.Closer
	router *routing.Router

	Database struct {
		MysqlDB *database.MysqlDB
	}

	Repository struct {
		UserRepository    repository.UserRepository
		FriendRepository  repository.FriendRepository
		ProfileRepository repository.ProfileRepository
	}

	Service struct {
		AuthService    service.AuthService
		UserService    service.UserService
		FriendService  service.FriendService
		ProfileService service.ProfileService
	}

	Usecases struct {
		SignInUsecase *usecase.SignInUsecase
		SignUpUsecase *usecase.SignUpUsecase

		GetProfileUsecase     *usecase.GetProfileUsecase
		SearchProfilesUsecase *usecase.SearchProfilesUsecase
		AddFriendUsecase      *usecase.AddFriendUsecase
		RemoveFriendUsecase   *usecase.RemoveFriendUsecase
	}
}

func New(cfgPath string) *DIContainer {
	cfg := config.NewConfig(cfgPath)

	dc := &DIContainer{
		Config: cfg,
		closer: &closer.Closer{},
	}

	dc.Database.MysqlDB = database.NewMsSQLDatabase(cfg)
	dc.closer.Add(func(ctx context.Context) error {
		return dc.Database.MysqlDB.Close()
	})

	dc.buildRepositories()
	dc.buildServices()
	dc.buildUsecases()
	dc.buildRouter()

	return dc
}

func (dc *DIContainer) Start() {
	dc.router.Start()
}

func (dc *DIContainer) Stop(ctx context.Context) {
	dc.closer.Close(ctx)
}

func (dc *DIContainer) buildRepositories() {
	dc.Repository.UserRepository = userRepository.NewMysqlUserRepository(
		dc.Database.MysqlDB,
	)

	dc.Repository.FriendRepository = friendRepository.NewMysqlFriendRepository(
		dc.Database.MysqlDB,
	)

	dc.Repository.ProfileRepository = profileRepository.NewMysqlProfileRepository(
		dc.Database.MysqlDB,
	)
}

func (dc *DIContainer) buildServices() {
	dc.Service.AuthService = service.NewAuthService(
		dc.Config,
	)

	dc.Service.UserService = service.NewUserService(
		dc.Repository.UserRepository,
	)

	dc.Service.FriendService = service.NewFriendService(
		dc.Repository.FriendRepository,
	)

	dc.Service.ProfileService = service.NewProfileService(
		dc.Repository.ProfileRepository,
	)
}

func (dc *DIContainer) buildUsecases() {
	dc.Usecases.SignInUsecase = usecase.NewSignInUsecase(
		dc.Service.AuthService,
		dc.Service.UserService,
	)

	dc.Usecases.SignUpUsecase = usecase.NewSignUpUsecase(
		dc.Service.AuthService,
		dc.Service.UserService,
	)

	dc.Usecases.SearchProfilesUsecase = usecase.NewSearchProfilesUsecase(
		dc.Service.ProfileService,
	)

	dc.Usecases.GetProfileUsecase = usecase.NewGetProfileUsecase(
		dc.Service.ProfileService,
	)

	dc.Usecases.AddFriendUsecase = usecase.NewAddFriendUsecase(
		dc.Service.FriendService,
	)

	dc.Usecases.RemoveFriendUsecase = usecase.NewRemoveFriendUsecase(
		dc.Service.FriendService,
	)
}

func (dc *DIContainer) buildRouter() {
	s, err := server.NewServer(nil)
	if err != nil {
		log.Fatalf("failed to init server %v", err)
	}

	// static
	s.AddTemplates("templates/**/*.tmpl")
	s.AddStatic("/assets", "templates/assets")

	r := routing.NewRouter(
		s,
		dc.Config.Http.Domain,
		dc.Service.AuthService,
		dc.Usecases.SignInUsecase,
		dc.Usecases.SignUpUsecase,
		dc.Usecases.GetProfileUsecase,
		dc.Usecases.SearchProfilesUsecase,
		dc.Usecases.AddFriendUsecase,
		dc.Usecases.RemoveFriendUsecase,
	)

	dc.router = r

	dc.closer.Add(r.Stop)
}
