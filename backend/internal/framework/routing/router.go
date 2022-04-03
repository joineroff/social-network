package routing

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joineroff/social-network/backend/internal/framework/routing/handler/api"
	"github.com/joineroff/social-network/backend/internal/framework/routing/handler/ssr"
	"github.com/joineroff/social-network/backend/internal/framework/routing/middleware"
	"github.com/joineroff/social-network/backend/internal/service"
	"github.com/joineroff/social-network/backend/internal/usecase"
	"github.com/joineroff/social-network/backend/pkg/server"
)

type Router struct {
	server         *server.Server
	authMiddleware *middleware.AuthMiddleware
	domain         string
}

func NewRouter(
	s *server.Server,
	domain string,
	authService service.AuthService,
	signInUsecase *usecase.SignInUsecase,
	signUpUsecase *usecase.SignUpUsecase,
	getProfileUsecase *usecase.GetProfileUsecase,
	searchProfilesUsecase *usecase.SearchProfilesUsecase,
	addFriendUsecase *usecase.AddFriendUsecase,
	removeFriendUsecase *usecase.RemoveFriendUsecase,
) *Router {
	r := &Router{
		server: s,
		domain: domain,
	}

	r.authMiddleware = middleware.NewAuthMiddleware(authService)

	// API
	s.AddRoute(http.MethodGet, "/api/v1/health", []gin.HandlerFunc{api.HealthCheckHandler()})
	s.AddRoute(http.MethodPost, "/api/v1/auth/sign-in", []gin.HandlerFunc{
		api.SignInHandler(signInUsecase),
	})
	s.AddRoute(http.MethodPost, "/api/v1/auth/sign-up", []gin.HandlerFunc{
		api.SignUpHandler(signUpUsecase),
	})
	s.AddRoute(http.MethodPost, "/api/v1/profile/search", []gin.HandlerFunc{
		api.SearchProfileHandler(searchProfilesUsecase),
	})

	// SSR
	s.AddRoute(http.MethodGet, "/", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.HomeHandler(r.domain, getProfileUsecase),
	})
	s.AddRoute(http.MethodGet, "/sign-in", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.SignInHandler(),
	})
	s.AddRoute(http.MethodPost, "/sign-in", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.SignInPostHandler(r.domain, signInUsecase),
	})
	s.AddRoute(http.MethodGet, "/sign-up", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.SignUpHandler(),
	})
	s.AddRoute(http.MethodPost, "/sign-up", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.SignUpPostHandler(r.domain, signUpUsecase),
	})
	s.AddRoute(http.MethodGet, "/sign-out", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.SignOutPostHandler(r.domain),
	})
	s.AddRoute(http.MethodPost, "/sign-out", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.SignOutPostHandler(r.domain),
	})
	s.AddRoute(http.MethodGet, "/search", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.SearchHandler(searchProfilesUsecase),
	})

	s.AddRoute(http.MethodGet, "/profile/:profileID", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.ProfileHandler(getProfileUsecase),
	})

	s.AddRoute(http.MethodPost, "/profile/:profileID/addFriend", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.ProfileAddFriendHandler(addFriendUsecase),
	})

	s.AddRoute(http.MethodPost, "/profile/:profileID/removeFriend", []gin.HandlerFunc{
		r.authMiddleware.ExtractUserID(),
		ssr.ProfileRemoveFriendHandler(removeFriendUsecase),
	})

	return r
}

func (r *Router) Start() {
	r.server.Start()
}

func (r *Router) Stop(ctx context.Context) error {
	return r.server.Stop(ctx)
}
