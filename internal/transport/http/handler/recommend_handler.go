package handler

import "net/http"

type RecommendHandler struct{}

func NewRecommend() *RecommendHandler { return &RecommendHandler{} }
func (h *RecommendHandler) Home(w http.ResponseWriter, req *http.Request) {}
