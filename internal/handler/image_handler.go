package handler

import (
	"encoding/json"
	"net/http"

	"Image-Hosting/internal/service"
)

type ImageHandler struct {
	svc *service.ImageService
}

func NewImageHandler(svc *service.ImageService) *ImageHandler {
	return &ImageHandler{svc: svc}
}

type createReq struct {
	Filename string `json:"filename"`
}

type createResp struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	PresignedURL string `json:"presigned_url"`
}

func (h *ImageHandler) CreateImage(w http.ResponseWriter, r *http.Request) {
	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Filename == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	id, url, err := h.svc.CreatePendingImage(req.Filename)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createResp{
		ID:           id,
		Status:       "PENDING",
		PresignedURL: url,
	})
}

func (h *ImageHandler) UploadToS3(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h.svc.ConfirmS3Upload(id)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "file uploaded to S3"}`))
}

func (h *ImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	img, cacheStatus, err := h.svc.GetMetadata(id)
	if err != nil {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("X-Cache", cacheStatus)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(img)
}
