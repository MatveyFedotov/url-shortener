package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"url-shortener/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

type CreateRequest struct {
	URL string `json:"url"`
}

type CreateResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "url is required",
		})
		return
	}

	code, err := h.service.CreateShortURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	shortURL := "http://localhost:8080/" + code

	c.JSON(http.StatusOK, CreateResponse{
		Code:     code,
		ShortURL: shortURL,
	})
}

func (h *Handler) Get(c *gin.Context) {
	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "code is required",
		})
		return
	}

	originalURL, err := h.service.GetOriginalURL(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "url not found",
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
