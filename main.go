package main

import (
	"context"
	"example/catalogue/config"
	"example/catalogue/data"
	"example/catalogue/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func main() {
	cfg := config.GetConfig()

	ctx := context.Background()
	exp, err := newExporter(cfg.JaegerAddress)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)
	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	// Finally, set the tracer that can be used for this package.
	tracer := tp.Tracer("catalogue-service")
	otel.SetTextMapPropagator(propagation.TraceContext{})

	productRepository := data.ProductRepository{Tracer: tracer, Products: data.Products}
	productHandler := handlers.ProductHandler{Tracer: tracer, Repo: productRepository}

	router := mux.NewRouter()
	router.HandleFunc("/{id}", productHandler.GetProduct).Methods("GET")

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM)

	<-stopCh

	if err := server.Close(); err != nil {
		log.Fatal("error while stopping server: ", err)
	}
}

func newExporter(address string) (*jaeger.Exporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(address)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("catalogue-service"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
