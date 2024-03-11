package handler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/christianferraz/tempwithopentelemetry/microservice/service-b/internal/entity"
	"github.com/go-chi/chi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func BuscaCepHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	tr := otel.Tracer("microservice-trace")

	data, cepOutput, err := BuscaCep(cep, tr)
	if err != nil {
		if err.Error() == "can not find zipcode" {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro ao buscar CEP", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entity.OutPutDTO{
		City:  cepOutput.Localidade,
		TempC: data.Current.TempC,
		TempF: data.Current.TempF,
		TempK: data.Current.TempC + 273.15,
	})
}

func BuscaCep(cep string, tr trace.Tracer) (*entity.Weather, *entity.CepDTOOutput, error) {
	_, span := tr.Start(context.Background(), "get cep")

	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, nil := io.ReadAll(resp.Body)
	if err != nil {
		return &entity.Weather{}, &entity.CepDTOOutput{}, err
	}

	// Se for um erro, deserializa para a estrutura de erro
	var erroResp entity.ErroCEP
	if json.Unmarshal(body, &erroResp) == nil && erroResp.Erro {

		return &entity.Weather{}, &entity.CepDTOOutput{}, fmt.Errorf("can not find zipcode")
	}
	var erroResp2 entity.ErroCEP2
	if json.Unmarshal(body, &erroResp2) == nil && erroResp2.Erro != "" {

		return &entity.Weather{}, &entity.CepDTOOutput{}, fmt.Errorf("can not find zipcode")
	}

	// Se não for um erro, deserializa para a estrutura de sucesso
	var cepResp entity.CepDTOOutput
	if err := json.Unmarshal(body, &cepResp); err != nil {
		log.Fatalf("Erro ao deserializar resposta de sucesso: %v", err)
	}
	span.End()
	weather, err := BuscaWeather(&cepResp, tr)

	return weather, &cepResp, err
}

func BuscaWeather(cep *entity.CepDTOOutput, tr trace.Tracer) (*entity.Weather, error) {
	_, span := tr.Start(context.Background(), "get cep")
	defer span.End()
	token := "a81637a90bbb4bce98a45909240603"
	req := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", token, url.QueryEscape(cep.Localidade))
	resp, err := http.Get(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data entity.Weather
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil

}
