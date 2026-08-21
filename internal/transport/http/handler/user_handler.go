package handler

import "net/http"

type UserHandler struct{}

func NewUser() *UserHandler { return &UserHandler{} }
func (h *UserHandler) Me(w http.ResponseWriter, req *http.Request)       {}
func (h *UserHandler) Password(w http.ResponseWriter, req *http.Request) {}
