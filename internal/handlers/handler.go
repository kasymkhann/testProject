package handlers

import (
	"net/http"
	"strconv"
	"testProject/internal/model"
	"testProject/pkg/logging"
	"testProject/service"

	"github.com/gin-gonic/gin"
)

// Handler представляет собой обработчик HTTP-запросов для взаимодействия с сервисом.
type Handler struct {
	service *service.Service
	logger  *logging.Logger
}

// NewHandler принимает service и logger в конструкторе и возрашает cтруктуру *Handler.
func NewHandler(service service.Service, logger *logging.Logger) *Handler {
	return &Handler{service: &service, logger: logger}
}

// CreatePerson обработчик создания нового человека.
func (h *Handler) CreatePerson(c *gin.Context) {
	h.logger.Debug("Handling CreatePerson request")
	var input model.Person
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Fatalf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if err := h.service.CreatePerson(&input); err != nil {
		h.logger.Errorf("Failed to create person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create person"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": input.ID})
}

// GetPeople обработчик получения списка людей.
func (h *Handler) GetPeople(c *gin.Context) {
	h.logger.Debug("Handling GetPeople request")
	filters := make(map[string]interface{})
	for key, value := range c.Request.URL.Query() {
		if len(value) > 0 {
			filters[key] = value[0]
		}
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		h.logger.Errorf("Failed to parse offset: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.logger.Errorf("Failed to parse limit: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}
	people, err := h.service.GetPeople(filters, offset, limit)
	if err != nil {
		h.logger.Errorf("Failed to get people: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get people"})
		return
	}
	c.JSON(http.StatusOK, people)
}

// GetPersonById обработчик получения информации о человеке по идентификатору.
func (h *Handler) GetPersonById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Errorf("Failed to parse person ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person ID"})
		return
	}

	persone, err := h.service.GetPersonById(id)
	if err != nil {
		h.logger.Errorf("Failed to get person by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get person by ID"})
		return
	}
	c.JSON(http.StatusOK, persone)

}

// UpdatePerson обработчик обновления информации о человеке.
func (h *Handler) UpdatePerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Fatalf("Failed to parse person ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person ID"})
		return
	}

	var input model.Person

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	input.ID = uint(id)

	if err := h.service.UpdatePerson(&input); err != nil {
		h.logger.Errorf("Failed to update person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update person"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "person updated successfully"})
}

// DeletePerson обработчик удаления информации о человеке.
func (h *Handler) DeletePerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Fatalf("Failed to parse person ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person ID"})
		return
	}

	if err := h.service.DeletePerson(id); err != nil {
		h.logger.Fatalf("Failed to delete person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete person"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "person deleted successfully"})

}
