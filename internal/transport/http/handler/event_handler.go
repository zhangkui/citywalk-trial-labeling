package handler

import "net/http"

type EventHandler struct{}

func NewEvent() *EventHandler { return &EventHandler{} }
func (h *EventHandler) List(w http.ResponseWriter, req *http.Request)       {}
func (h *EventHandler) Detail(w http.ResponseWriter, req *http.Request)     {}
func (h *EventHandler) Create(w http.ResponseWriter, req *http.Request)     {}
func (h *EventHandler) Update(w http.ResponseWriter, req *http.Request)     {}
func (h *EventHandler) Delete(w http.ResponseWriter, req *http.Request)     {}
func (h *EventHandler) Join(w http.ResponseWriter, req *http.Request)       {}
func (h *EventHandler) CancelJoin(w http.ResponseWriter, req *http.Request) {}
func (h *EventHandler) CheckIn(w http.ResponseWriter, req *http.Request)    {}
