package repository

import (
	"context"

	"github.com/thiiagofernando/cep-geo-api/internal/domain/entity"
	"github.com/thiiagofernando/cep-geo-api/internal/domain/repository"
	"github.com/thiiagofernando/cep-geo-api/internal/infra/geocoder"
)

type AddressRepositoryImpl struct {
	geocoder *geocoder.Geocoder
}

func NewAddressRepository() repository.AddressRepository {
	return &AddressRepositoryImpl{
		geocoder: geocoder.GetGeocoder(),
	}
}

func (r *AddressRepositoryImpl) GetAddressByCEP(ctx context.Context, cep string) (*entity.Address, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Continua o processamento
	}

	return r.geocoder.GeocodeAddress(cep)
}
