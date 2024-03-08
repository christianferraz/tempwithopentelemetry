package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials/insecure"
)

type CEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ErroCEP struct {
	Erro bool `json:"erro"`
}

type ErroCEP2 struct {
	Erro string `json:"erro"`
}

type Weather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdated string  `json:"last_updated"`
		TempC       float64 `json:"temp_c"`
		TempF       float64 `json:"temp_f"`
		IsDay       int     `json:"is_day"`
		Condition   struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

type OutPutDTO struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func initProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create resource: %w",
			err,
		)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create gRPC connection to collector: %w",
			err,
		)
	}
	defer conn.Close()
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	return tracerProvider.Shutdown, nil
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /cep/{id}", BuscaCepHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func BuscaCepHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("id")
	pathSegments := strings.Split(r.URL.Path, "/")
	if cep == "" {
		cep = pathSegments[len(pathSegments)-1]
	}
	if !cepFormatted(cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}
	data, err := BuscaCep(cep)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(OutPutDTO{
		TempC: data.Current.TempC,
		TempF: data.Current.TempF,
		TempK: data.Current.TempC + 273.15,
	})
}

func BuscaCep(cep string) (*Weather, error) {
	transport := http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Se for um erro, deserializa para a estrutura de erro
	var erroResp ErroCEP
	if json.Unmarshal(body, &erroResp) == nil && erroResp.Erro {

		return nil, fmt.Errorf("can not find zipcode")
	}
	var erroResp2 ErroCEP2
	if json.Unmarshal(body, &erroResp2) == nil && erroResp2.Erro != "" {

		return nil, fmt.Errorf("can not find zipcode")
	}

	// Se não for um erro, deserializa para a estrutura de sucesso
	var cepResp CEP
	if err := json.Unmarshal(body, &cepResp); err != nil {
		log.Fatalf("Erro ao deserializar resposta de sucesso: %v", err)
	}

	return BuscaWeather(&cepResp)
}

func BuscaWeather(cep *CEP) (*Weather, error) {
	token := "a81637a90bbb4bce98a45909240603"
	req := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", token, url.QueryEscape(cep.Localidade))
	resp, err := http.Get(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data Weather
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil

}

func cepFormatted(input string) bool {
	matched, err := regexp.MatchString(`^[0-9]{8}$`, input)
	if err != nil {
		fmt.Printf("Erro ao verificar a string: %v\n", err)
		return false
	}
	return matched
}
