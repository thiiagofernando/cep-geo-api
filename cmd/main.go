package main

import (
	"fmt"
	"log"

	"github.com/thiiagofernando/cep-geo-api/internal/infra/api"
	"github.com/thiiagofernando/cep-geo-api/internal/infra/api/handler"
	"github.com/thiiagofernando/cep-geo-api/internal/infra/config"
	repository "github.com/thiiagofernando/cep-geo-api/internal/infra/repository"
	"github.com/thiiagofernando/cep-geo-api/internal/usecase"
)

// @title CEP Geolocalização API
// @version 1.0
// @description API para obter coordenadas geográficas (latitude e longitude) a partir de um CEP brasileiro
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8070
// @BasePath /api/v1
func main() {
	// Carrega configurações
	cfg := config.LoadConfig()

	// Inicializa o repositório
	addressRepo := repository.NewAddressRepository()

	// Inicializa o caso de uso
	getAddressByCEPUseCase := usecase.NewGetAddressByCEPUseCase(addressRepo)

	// Inicializa o handler
	addressHandler := handler.NewAddressHandler(getAddressByCEPUseCase)

	// Configura o router
	router := api.SetupRouter(addressHandler)

	// Inicia o servidor
	log.Printf("Servidor iniciado na porta %d", cfg.Port)
	log.Printf("Documentação Swagger disponível em http://localhost:%d/swagger/index.html", cfg.Port)

	if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
