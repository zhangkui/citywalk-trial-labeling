package handler

import "net/http"

type SocialHandler struct{}

func NewSocial() *SocialHandler { return &SocialHandler{} }
func (h *SocialHandler) Comments(w http.ResponseWriter, req *http.Request)     {}
func (h *SocialHandler) CreateComment(w http.ResponseWriter, req *http.Request) {}
func (h *SocialHandler) UpdateComment(w http.ResponseWriter, req *http.Request) {}
func (h *SocialHandler) DeleteComment(w http.ResponseWriter, req *http.Request) {}
