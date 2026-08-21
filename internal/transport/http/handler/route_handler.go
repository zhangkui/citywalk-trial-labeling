package handler

import "net/http"

type RouteHandler struct{}

func NewRoute() *RouteHandler { return &RouteHandler{} }
func (h *RouteHandler) List(w http.ResponseWriter, req *http.Request)       {}
func (h *RouteHandler) Detail(w http.ResponseWriter, req *http.Request)     {}
func (h *RouteHandler) Create(w http.ResponseWriter, req *http.Request)     {}
func (h *RouteHandler) Update(w http.ResponseWriter, req *http.Request)     {}
func (h *RouteHandler) Delete(w http.ResponseWriter, req *http.Request)     {}
func (h *RouteHandler) Status(w http.ResponseWriter, req *http.Request)     {}
func (h *RouteHandler) Favorite(w http.ResponseWriter, req *http.Request)   {}
func (h *RouteHandler) Unfavorite(w http.ResponseWriter, req *http.Request) {}
func (h *RouteHandler) Rate(w http.ResponseWriter, req *http.Request)       {}
