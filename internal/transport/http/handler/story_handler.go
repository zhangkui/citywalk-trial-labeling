package handler

import "net/http"

type StoryHandler struct{}

func NewStory() *StoryHandler { return &StoryHandler{} }
func (h *StoryHandler) List(w http.ResponseWriter, req *http.Request)   {}
func (h *StoryHandler) Detail(w http.ResponseWriter, req *http.Request) {}
func (h *StoryHandler) Create(w http.ResponseWriter, req *http.Request) {}
func (h *StoryHandler) Update(w http.ResponseWriter, req *http.Request) {}
func (h *StoryHandler) Delete(w http.ResponseWriter, req *http.Request) {}
func (h *StoryHandler) Status(w http.ResponseWriter, req *http.Request) {}
func (h *StoryHandler) Like(w http.ResponseWriter, req *http.Request)   {}
func (h *StoryHandler) Unlike(w http.ResponseWriter, req *http.Request) {}
