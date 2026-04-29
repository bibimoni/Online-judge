package server

import (
	"net"
	"net/http"

	"github.com/bibimoni/Online-judge/gateway/src/infrastructure/config"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
}

var server *Server = nil

func NewServer() *Server {
	cfg := config.Load()
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	if server == nil {
		r := chi.NewRouter()
		server = &Server{
			router: r,
			httpServer: &http.Server{
				Addr:    serverAddr,
				Handler: r,
			},
		}
	}
	return server
}

func (s *Server) ListenAndServe() {
	l, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		panic(err)
	}

	log := config.GetLogger()
	log.Info().Msgf("API gateway is listening on: %s", s.httpServer.Addr)

	s.httpServer.Serve(l)
}

func GetRouter() *chi.Mux {
	return server.router
}

func GetServer() *http.Server {
	return server.httpServer
}
