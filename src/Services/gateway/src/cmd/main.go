package main

import (
	"github.com/bibimoni/Online-judge/gateway/src/infrastructure/config"
	"github.com/bibimoni/Online-judge/gateway/src/middlewares"
	"github.com/bibimoni/Online-judge/gateway/src/proxy"
	"github.com/bibimoni/Online-judge/gateway/src/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	config.NewLogger(cfg.LogLevel)
	s := server.NewServer()
	r := server.GetRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", proxy.LoginApiProxy())

	r.Route("/api/v1/submission", func(r chi.Router) {
		r.Method("GET", "/view/*", proxy.SubmissionApiProxy())
		r.Method("GET", "/problem/view/*", proxy.SubmissionApiProxy())
		r.With(middlewares.WithAuth).Method("POST", "/submit", proxy.SubmissionApiProxy())
		r.Handle("/ws", proxy.WSSubmissionProxy(cfg.Endpoints.Submission))
	})

	r.Route("/problem", func(r chi.Router) {
		r.Method("GET", "/all", proxy.ProblemApiProxy())
		r.Method("GET", "/get/{param}/statement.pdf", proxy.ProblemApiProxy())
		r.Method("GET", "/get/{param}/problem.json", proxy.ProblemApiProxy())
	})

	s.ListenAndServe()
}
