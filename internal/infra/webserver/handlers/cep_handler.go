package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aleroxac/goexpert-fatest-zipcode-api/internal/dto"
	"github.com/aleroxac/goexpert-fatest-zipcode-api/internal/entity"
	"github.com/go-chi/chi"
)

type CEPHandler struct{}

func NewGetCEPHandler() *CEPHandler {
	return &CEPHandler{}
}

func GetCEP(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatalf("Fail to create the request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Fail to make the request: %v", err)
	}
	defer res.Body.Close()

	ctx_err := ctx.Err()
	if ctx_err != nil {
		select {
		case <-ctx.Done():
			err := ctx.Err()
			log.Fatalf("Max timeout reached: %v", err)
		}
	}

	resp_json, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Fail to read the response: %v", err)
	}

	return resp_json, nil
}

func GetCEPBrasilAPI(cep string) (entity.BrasilApiOutput, error) {
	api_response_body, err := GetCEP(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep))
	if err != nil {
		panic(err)
	}

	var cep_unmarshaled entity.BrasilApiOutput
	err = json.Unmarshal(api_response_body, &cep_unmarshaled)
	if err != nil {
		panic(err)
	}

	return cep_unmarshaled, nil
}

func GetViaCEPAPI(cep string) (entity.ViaCepOutput, error) {
	api_response_body, err := GetCEP(fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep))
	if err != nil {
		panic(err)
	}

	var cep_unmarshaled entity.ViaCepOutput
	err = json.Unmarshal(api_response_body, &cep_unmarshaled)
	if err != nil {
		panic(err)
	}

	return cep_unmarshaled, nil
}

// GetCEP godoc
//
//	@Summary		Get CEP
//	@Description	Get CEP
//	@Tags			cep
//	@Accept			json
//	@Produce		json
//	@Param			cep				path		string	true	"cep address"
//	@Success		200				{object}	interface{}
//	@Failure		500				{object}	dto.Error
//	@Router			/cep/{cep}    	[get]
func (h *CEPHandler) GetCEP(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	brasilapi_chan := make(chan entity.BrasilApiOutput)
	viacep_chan := make(chan entity.ViaCepOutput)

	// brasilapi
	go func() {
		// time.Sleep(time.Second * 2)
		msg, err := GetCEPBrasilAPI(cep)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			error := dto.Error{Message: err.Error()}
			json.NewEncoder(w).Encode(error)
			return
		}
		brasilapi_chan <- msg
	}()

	// viacep
	go func() {
		msg, err := GetViaCEPAPI(cep)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			error := dto.Error{Message: err.Error()}
			json.NewEncoder(w).Encode(error)
			return
		}
		viacep_chan <- msg
	}()

	select {
	case brasilapi := <-brasilapi_chan:
		cep_data, err := json.Marshal(brasilapi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			error := dto.Error{Message: err.Error()}
			json.NewEncoder(w).Encode(error)
			return
		}
		fmt.Println(string(cep_data))

		cep_output := dto.CEPOutput{
			FatestAPI: "brasilapi",
			Response:  brasilapi,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cep_output)
	case viacep := <-viacep_chan:
		cep_data, err := json.Marshal(viacep)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			error := dto.Error{Message: err.Error()}
			json.NewEncoder(w).Encode(error)
			return
		}
		fmt.Println("viacep:", string(cep_data))

		cep_output := dto.CEPOutput{
			FatestAPI: "viacep",
			Response:  viacep,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cep_output)
	case <-time.After(time.Second * 1):
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("timeout")
		fmt.Println("timeout")
	}
}
