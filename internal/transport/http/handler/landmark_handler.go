package handler

import "net/http"

type LandmarkHandler struct{}

func NewLandmark() *LandmarkHandler { return &LandmarkHandler{} }
func (h *LandmarkHandler) List(w http.ResponseWriter, req *http.Request)   {}
func (h *LandmarkHandler) Detail(w http.ResponseWriter, req *http.Request) {}
func (h *LandmarkHandler) Create(w http.ResponseWriter, req *http.Request) {}
func (h *LandmarkHandler) Update(w http.ResponseWriter, req *http.Request) {}
func (h *LandmarkHandler) Delete(w http.ResponseWriter, req *http.Request) {}
func (h *LandmarkHandler) Import(w http.ResponseWriter, req *http.Request) {}
func (h *LandmarkHandler) Nearby(w http.ResponseWriter, req *http.Request) {}
