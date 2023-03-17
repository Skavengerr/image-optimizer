package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Skavengerr/image-optimizer/model"
	"github.com/Skavengerr/image-optimizer/queue"
	"github.com/Skavengerr/image-optimizer/service"
	"github.com/Skavengerr/image-optimizer/util"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ImageHandler struct {
	service      *service.ImageService
	messageQueue queue.MessageQueue
}

type ImageIDResponse struct {
	ID string `json:"id"`
}

func NewImageHandler(service *service.ImageService, messageQueue queue.MessageQueue) *ImageHandler {
	return &ImageHandler{
		service:      service,
		messageQueue: messageQueue,
	}
}

func (h *ImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Missing or invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageID := uuid.New().String()

	f, err := os.OpenFile(fmt.Sprintf("./images/%s.jpg", imageID), os.O_WRONLY|os.O_CREATE, 0666)
	fmt.Println("err", err)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	if err := h.sendMessageToQueue(imageID); err != nil {
		http.Error(w, "Failed to queue image for processing", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": imageID})
}

func (h *ImageHandler) DownloadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	qualityStr := r.URL.Query().Get("quality")
	quality, err := strconv.Atoi(qualityStr)
	if err != nil || quality < 0 || quality > 100 {
		http.Error(w, "Invalid quality parameter", http.StatusBadRequest)
		return
	}
	scale := float64(quality) / 100.0

	imagePath := fmt.Sprintf("./images/%s.jpg", id)
	resizedImg, err := util.ResizeImage(imagePath, scale)
	fmt.Println("err4", err)
	if err != nil {
		http.Error(w, "Failed to resize image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")

	_, err = w.Write(resizedImg)
	if err != nil {
		http.Error(w, "Failed to write image to response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ImageHandler) sendMessageToQueue(imageID string) error {
	message := &model.Message{
		ID: imageID,
	}

	err := h.messageQueue.SendMessage(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send message to RabbitMQ queue: %w", err)
	}

	return nil
}
