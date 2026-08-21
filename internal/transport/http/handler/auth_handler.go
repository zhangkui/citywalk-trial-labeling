package handler

import "net/http"

type AuthHandler struct{}

func NewAuth() *AuthHandler { return &AuthHandler{} }

func (h *AuthHandler) Register(w http.ResponseWriter, req *http.Request)  {}
func (h *AuthHandler) Login(w http.ResponseWriter, req *http.Request)     {}
func (h *AuthHandler) Refresh(w http.ResponseWriter, req *http.Request)   {}
func (h *AuthHandler) Logout(w http.ResponseWriter, req *http.Request)     {}
func (h *AuthHandler) Me(w http.ResponseWriter, req *http.Request)        {}
func (h *AuthHandler) Password(w http.ResponseWriter, req *http.Request)  {}
