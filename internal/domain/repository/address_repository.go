package repository

import (
	"context"

	"github.com/thiiagofernando/cep-geo-api/internal/domain/entity"
)

type AddressRepository interface {
	GetAddressByCEP(ctx context.Context, cep string) (*entity.Address, error)
}
