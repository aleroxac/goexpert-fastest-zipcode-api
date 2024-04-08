package dto

type CEPInput struct {
	CEP string `json:"cep"`
}

type CEPOutput struct {
	FatestAPI string      `json:"fatest_api"`
	Response  interface{} `json:"response"`
}

type Error struct {
	Message string `json:"message"`
}
