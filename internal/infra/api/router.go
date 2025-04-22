package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/thiiagofernando/cep-geo-api/docs" 
	"github.com/thiiagofernando/cep-geo-api/internal/infra/api/handler"
	"github.com/thiiagofernando/cep-geo-api/internal/infra/api/middleware"
)


func SetupRouter(addressHandler *handler.AddressHandler) *gin.Engine {
	r := gin.Default()

	// Adiciona middlewares globais
	r.Use(middleware.CORS())

	// Configura grupo de rotas para API
	api := r.Group("/api/v1")
	{
		api.GET("/cep/:cep", addressHandler.GetAddressByCEP)
	}

	// Configura Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
