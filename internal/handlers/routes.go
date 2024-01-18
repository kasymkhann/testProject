package handlers

import (
	"testProject/pkg/logging"
	"testProject/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes регистрирует маршруты HTTP для взаимодействия с обработчиками, используемыми сервисом.
func RegisterRoutes(router *gin.Engine, service *service.Service) {
	handler := NewHandler(*service, logging.GetLogger())

	router.POST("/people", handler.CreatePerson)
	router.GET("/people", handler.GetPeople)
	router.GET("/people/:id", handler.GetPersonById)
	router.PUT("/people/:id", handler.UpdatePerson)
	router.DELETE("/people/:id", handler.DeletePerson)

}
