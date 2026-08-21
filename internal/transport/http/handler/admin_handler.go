package handler

import "net/http"

type AdminHandler struct{}

func NewAdmin() *AdminHandler { return &AdminHandler{} }
func (h *AdminHandler) Users(w http.ResponseWriter, req *http.Request)         {}
func (h *AdminHandler) CreateUser(w http.ResponseWriter, req *http.Request)    {}
func (h *AdminHandler) UserStatus(w http.ResponseWriter, req *http.Request)    {}
func (h *AdminHandler) UserPassword(w http.ResponseWriter, req *http.Request)  {}
func (h *AdminHandler) Roles(w http.ResponseWriter, req *http.Request)         {}
func (h *AdminHandler) CreateRole(w http.ResponseWriter, req *http.Request)    {}
func (h *AdminHandler) UpdateRole(w http.ResponseWriter, req *http.Request)    {}
func (h *AdminHandler) DeleteRole(w http.ResponseWriter, req *http.Request)    {}
func (h *AdminHandler) Permissions(w http.ResponseWriter, req *http.Request)   {}
func (h *AdminHandler) AssignRoles(w http.ResponseWriter, req *http.Request)   {}
