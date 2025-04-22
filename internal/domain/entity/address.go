package entity

type Address struct {
	CEP       string  `json:"cep"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Street    string  `json:"street"`
	City      string  `json:"city"`
	State     string  `json:"state"`
}
