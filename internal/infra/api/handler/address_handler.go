package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thiiagofernando/cep-geo-api/internal/usecase"
)

type AddressHandler struct {
	getAddressByCEPUseCase *usecase.GetAddressByCEPUseCase
}

func NewAddressHandler(getAddressByCEPUseCase *usecase.GetAddressByCEPUseCase) *AddressHandler {
	return &AddressHandler{
		getAddressByCEPUseCase: getAddressByCEPUseCase,
	}
}

// GetAddressByCEP godoc
// @Summary Obtém coordenadas geográficas a partir de um CEP
// @Description Retorna a latitude e longitude de um endereço a partir do CEP informado
// @Tags endereços
// @Accept json
// @Produce json
// @Param cep path string true "CEP (apenas números ou formato 00000-000)"
// @Success 200 {object} entity.Address
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cep/{cep} [get]
func (h *AddressHandler) GetAddressByCEP(c *gin.Context) {
	cep := c.Param("cep")

	address, err := h.getAddressByCEPUseCase.Execute(c.Request.Context(), cep)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, address)
}
