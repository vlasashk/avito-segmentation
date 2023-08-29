package controller

import "net/http"

type APIServer struct {
	ListenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		ListenAddr: listenAddr,
	}
}

func (s *APIServer) HandleAccount(w http.ResponseWriter, r http.Request) error {
	return nil
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r http.Request) error {
	return nil
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r http.Request) error {
	return nil
}

func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r http.Request) error {
	return nil
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r http.Request) error {
	return nil
}
