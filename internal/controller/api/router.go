package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vlasashk/avito-segmentation/internal/model/storage"
	"log/slog"
	"net/http"
	"time"
)

type ServerAPI struct {
	ListenAddr string
	Store      *storage.PostgresDB
	Log        *slog.Logger
}

func NewAPIServer(listenAddr string, store *storage.PostgresDB, log *slog.Logger) *ServerAPI {
	return &ServerAPI{
		ListenAddr: listenAddr,
		Store:      store,
		Log:        log,
	}
}

func Run(log *slog.Logger, server *ServerAPI) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Use(middleware.Timeout(60 * time.Second))

	initRouterMethods()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello GolinuxCloud members!"))
	})

	router.Post("/user", server.HandleAddUser)
	router.Post("/segment", server.HandleAddSegment)

	if err := http.ListenAndServe(server.ListenAddr, router); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server died")
}

func initRouterMethods() {
	chi.RegisterMethod("GET")
	chi.RegisterMethod("POST")
	chi.RegisterMethod("PUT")
	chi.RegisterMethod("DELETE")
}
