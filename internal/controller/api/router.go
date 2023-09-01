package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

func NewAPIServer(listenAddr string, store Storage, log *slog.Logger) *ServerAPI {
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

	//initRouterMethods()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("Hello World!"))
	})
	router.Mount("/user", server.userRouter())
	router.Mount("/segment", server.segmentRouter())

	if err := http.ListenAndServe(":"+server.ListenAddr, router); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server died")
}

func (s *ServerAPI) userRouter() http.Handler {
	router := chi.NewRouter()
	router.Post("/new", s.HandleAddUser)
	router.Post("/addSegment", s.HandleAddUserToSegment)
	router.Get("/segments", s.HandleGetUserSegmentsInfo)
	router.Delete("/segments", s.HandleDeleteUserFromSegment)
	return router
}

func (s *ServerAPI) segmentRouter() http.Handler {
	router := chi.NewRouter()
	router.Post("/new", s.HandleAddSegment)
	return router
}

//func initRouterMethods() {
//	chi.RegisterMethod("GET")
//	chi.RegisterMethod("POST")
//	chi.RegisterMethod("PUT")
//}
