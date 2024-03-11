package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/christianferraz/tempwithopentelemetry/microservice/service-a/internal/entity"
	"github.com/christianferraz/tempwithopentelemetry/microservice/service-a/pkg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func BuscaCepHandler(w http.ResponseWriter, r *http.Request) {
	var httpclient http.Client
	tr := otel.Tracer("microservice-trace")
	ctx := context.Background()
	ctx, span := tr.Start(ctx, "get weather from service b")
	defer span.End()
	var cep entity.CepDTOInput
	err := json.NewDecoder(r.Body).Decode(&cep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if cep.Cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("zipcode is required"))
	}

	if !pkg.CepFormatted(&cep.Cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}
	url := fmt.Sprintf("%s/%s", "http://service-b:8081/cep", cep.Cep)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	request.Header.Set("Accept", "application/json")

	if err != nil {
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))
	response, err := httpclient.Do(request)
	if err != nil {
		fmt.Println(response.StatusCode)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(response.StatusCode)
		b, err := io.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(b))
		return
	}

	var weatherOutput entity.OutPutDTO
	err = json.NewDecoder(response.Body).Decode(&weatherOutput)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weatherOutput)
}
