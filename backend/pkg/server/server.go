package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Debug        bool          `envconfig:"DEBUG"`
	Host         string        `envconfig:"HOST"`
	Port         int           `envconfig:"PORT"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT"`
}

func defaultConfig() *ServerConfig {
	return &ServerConfig{
		Debug:        true,
		Host:         "0.0.0.0",
		Port:         80,
		ReadTimeout:  time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
}

type Server struct {
	engine *gin.Engine
	server *http.Server
}

type Route struct {
	Method      string
	Path        string
	HandleFuncs []gin.HandlerFunc
}

func NewServer(
	cfg *ServerConfig,
) (*Server, error) {
	engine := gin.New()

	if cfg == nil {
		cfg = defaultConfig()
	}

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{}

	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	s.engine = engine

	srv := &http.Server{
		Addr:         net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	s.server = srv

	return s, nil
}

func (s *Server) AddTemplates(tmplGlob string) {
	s.engine.LoadHTMLGlob(tmplGlob)
}

func (s *Server) AddStatic(relPath, root string) {
	s.engine.Static(relPath, root)
}

func (s *Server) AddRoute(
	method string,
	path string,
	handlers []gin.HandlerFunc,
) {
	s.engine.Handle(method, path, handlers...)
}

func (s *Server) Start() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("listen failed", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Serve will serve asynchronously (in a separate goroutine)
func (s *Server) Serve(ctx context.Context) {
	go func() {
		sctx, cancel := context.WithCancel(ctx)
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGQUIT)

		go func() {
			<-stopChan
			cancel()
		}()

		s.Start()

		defer func() {
			_ = s.Stop(context.Background())
		}()
		<-sctx.Done()
	}()
}
