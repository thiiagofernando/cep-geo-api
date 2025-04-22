package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/thiiagofernando/cep-geo-api/internal/domain/entity"
	"github.com/thiiagofernando/cep-geo-api/internal/domain/repository"
)

type GetAddressByCEPUseCase struct {
	repo repository.AddressRepository
}

func NewGetAddressByCEPUseCase(repo repository.AddressRepository) *GetAddressByCEPUseCase {
	return &GetAddressByCEPUseCase{
		repo: repo,
	}
}

func (uc *GetAddressByCEPUseCase) Execute(ctx context.Context, cep string) (*entity.Address, error) {

	cleanCEP := strings.ReplaceAll(cep, "-", "")
	cleanCEP = strings.ReplaceAll(cleanCEP, ".", "")
	cleanCEP = strings.TrimSpace(cleanCEP)

	matched, err := regexp.MatchString(`^\d{8}$`, cleanCEP)
	if err != nil || !matched {
		return nil, errors.New("formato de CEP inválido. Use 8 dígitos (00000000)")
	}

	return uc.repo.GetAddressByCEP(ctx, cleanCEP)
}
