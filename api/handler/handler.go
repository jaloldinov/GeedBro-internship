package handler

import (
	"auth/pkg/logger"
	"auth/storage"
)

type Handler struct {
	storage storage.StorageI
	log     logger.LoggerI
}

func NewHandler(strg storage.StorageI, loger logger.LoggerI) *Handler {
	return &Handler{storage: strg, log: loger}
}
