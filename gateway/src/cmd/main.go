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
		// r.With(middlewares.WithPermission("view_submission")).Method("GET", "/view/*", proxy.SubmissionApiProxy())
		// r.With(middlewares.WithPermission("view_problem")).Method("GET", "/problem/view/*", proxy.SubmissionApiProxy())
		r.Method("GET", "/view/*", proxy.SubmissionApiProxy())
		r.Method("GET", "/problem/view/*", proxy.SubmissionApiProxy())
		r.With(middlewares.WithPermission("submit_code")).Method("POST", "/submit", proxy.SubmissionApiProxy())
		r.With(middlewares.WithPermission("view_submission")).Handle("/ws", proxy.WSSubmissionProxy())
		r.Method("GET", "/lang/all", proxy.SubmissionApiProxy())
	})

	r.Route("/problem", func(r chi.Router) {
		r.Method("GET", "/all", proxy.ProblemApiProxy())
		r.Method("GET", "/get/{param}/statement.pdf", proxy.ProblemApiProxy())
		r.Method("GET", "/get/{param}/problem.json", proxy.ProblemApiProxy())
	})

	r.Route("/api/v1/contest", func(r chi.Router) {

		r.With(middlewares.OptionalAuth).Method("GET", "/", proxy.ContestApiProxy())
		r.With(middlewares.WithPermission("create_contest")).Method("POST", "/create", proxy.ContestApiProxy())
		r.With(middlewares.OptionalAuth).Method("POST", "/scoreboard", proxy.ContestApiProxy())

		r.Route("/{contest_id}", func(r chi.Router) {
			r.With(middlewares.OptionalAuth).Method("GET", "/", proxy.ContestApiProxy())
			r.With(middlewares.WithPermission("manage_contest")).Method("POST", "/edit", proxy.ContestApiProxy())
			r.With(middlewares.WithPermission("edit_contest")).Method("PATCH", "/patch", proxy.ContestApiProxy())
			r.With(middlewares.WithPermission("edit_contest")).Method("PUT", "/manage/problems", proxy.ContestApiProxy())

		})

		r.With(middlewares.WithPermission("register_contest")).Method("POST", "/register", proxy.ContestApiProxy())
		r.With(middlewares.WithPermission("register_contest")).Method("POST", "/unregister", proxy.ContestApiProxy())
		r.With(middlewares.WithPermission("submit_code")).Method("POST", "/submit", proxy.ContestApiProxy())
	})

	s.ListenAndServe()
}
