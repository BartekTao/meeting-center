package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	cloudstorage "cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type StorageHandler interface {
	Close() error
	UploadFile(file io.Reader, fileName string) (url string, err error)
}

type GCSHandler struct {
	client     *cloudstorage.Client
	bucketName string
}

func NewStorageHandler(client *cloudstorage.Client) StorageHandler {
	return &GCSHandler{client: client, bucketName: "nycu_meeting_center_1"}
}

func (h *GCSHandler) Close() error {
	return h.client.Close()
}

func (h *GCSHandler) UploadFile(file io.Reader, fileName string) (url string, err error) {
	ctx := context.Background()
	ext := filepath.Ext(fileName)
	objectName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	bucket := h.client.Bucket(h.bucketName)
	obj := bucket.Object(objectName)

	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write to GCS: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %v", err)
	}

	url = fmt.Sprintf("https://storage.googleapis.com/%s/%s", h.bucketName, objectName)
	return url, nil
}
