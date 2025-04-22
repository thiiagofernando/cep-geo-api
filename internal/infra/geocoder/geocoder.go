package geocoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/thiiagofernando/cep-geo-api/internal/domain/entity"
)


var (
	geocoderInstance *Geocoder
	once             sync.Once
)

type Geocoder struct {
	httpClient *http.Client
}


func GetGeocoder() *Geocoder {
	once.Do(func() {
		geocoderInstance = &Geocoder{
			httpClient: &http.Client{
				Timeout: 10 * time.Second,
			},
		}
	})
	return geocoderInstance
}


type ViaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
	Erro        bool   `json:"erro"`
}


type NominatimResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}


func (g *Geocoder) GeocodeAddress(cep string) (*entity.Address, error) {
	// Limpa o formato do CEP (remove traços, pontos, etc.)
	cleanCEP := strings.ReplaceAll(cep, "-", "")
	cleanCEP = strings.ReplaceAll(cleanCEP, ".", "")
	cleanCEP = strings.TrimSpace(cleanCEP)


	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cleanCEP)
	resp, err := g.httpClient.Get(viaCEPURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao consultar ViaCEP, status: %d", resp.StatusCode)
	}

	var viaCepResp ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCepResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta do ViaCEP: %w", err)
	}

	if viaCepResp.Erro {
		return nil, errors.New("CEP não encontrado")
	}

	// Adicionar tratamento para CEP sem logradouro
	if viaCepResp.Logradouro == "" {
		viaCepResp.Logradouro = "Centro" // Usa "Centro" como padrão quando não tiver logradouro
	}

	// Formatar o query para o Nominatim
	query := fmt.Sprintf("%s, %s, %s, Brasil", viaCepResp.Logradouro, viaCepResp.Localidade, viaCepResp.UF)
	
	// Codificar os parâmetros da URL para evitar problemas com caracteres especiais
	encodedQuery := url.QueryEscape(query)
	nominatimURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", encodedQuery)

	// Configurar a requisição com o User-Agent 
	req, err := http.NewRequest("GET", nominatimURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição para Nominatim: %w", err)
	}
	
	// Adicionar User-Agent  (obrigatório para Nominatim)
	req.Header.Add("User-Agent", "CEP-Geo-API/1.0 (github.com/thiiagofernando/cep-geo-api)")
	
	// Adicionar delay para respeitar o rate limiting do Nominatim
	time.Sleep(1 * time.Second)

	nominatimResp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar Nominatim: %w", err)
	}
	defer nominatimResp.Body.Close()

	if nominatimResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao consultar Nominatim, status: %d", nominatimResp.StatusCode)
	}

	var nominatimResults []NominatimResponse
	if err := json.NewDecoder(nominatimResp.Body).Decode(&nominatimResults); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta do Nominatim: %w", err)
	}

	if len(nominatimResults) == 0 {
		
		retryQuery := fmt.Sprintf("%s, %s, Brasil", viaCepResp.Localidade, viaCepResp.UF)
		encodedRetryQuery := url.QueryEscape(retryQuery)
		retryURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", encodedRetryQuery)
		
		retryReq, err := http.NewRequest("GET", retryURL, nil)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar requisição de retry para Nominatim: %w", err)
		}
		
		retryReq.Header.Add("User-Agent", "CEP-Geo-API/1.0 (github.com/thiiagofernando/cep-geo-api)")
		

		time.Sleep(1 * time.Second)
		
		retryResp, err := g.httpClient.Do(retryReq)
		if err != nil {
			return nil, fmt.Errorf("erro ao consultar Nominatim (retry): %w", err)
		}
		defer retryResp.Body.Close()
		
		if retryResp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("erro ao consultar Nominatim (retry), status: %d", retryResp.StatusCode)
		}
		
		if err := json.NewDecoder(retryResp.Body).Decode(&nominatimResults); err != nil {
			return nil, fmt.Errorf("erro ao decodificar resposta do Nominatim (retry): %w", err)
		}
		
		if len(nominatimResults) == 0 {
			return nil, errors.New("não foi possível obter as coordenadas para o CEP informado")
		}
	}

	// Parse das coordenadas
	var lat, lon float64
	fmt.Sscanf(nominatimResults[0].Lat, "%f", &lat)
	fmt.Sscanf(nominatimResults[0].Lon, "%f", &lon)

	return &entity.Address{
		CEP:       viaCepResp.CEP,
		Latitude:  lat,
		Longitude: lon,
		Street:    viaCepResp.Logradouro,
		City:      viaCepResp.Localidade,
		State:     viaCepResp.UF,
	}, nil
}
