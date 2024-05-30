package handlers

import (
	"encoding/json"
	"example/catalogue/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ProductHandler struct {
	Tracer trace.Tracer
	Repo   data.ProductRepository
}

func (h ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.Tracer.Start(r.Context(), "ProductHandler.GetProduct")
	defer span.End()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.Repo.GetProduct(ctx, int64(id))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResp, err := json.Marshal(product)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	span.SetStatus(codes.Ok, "")
	w.Write(jsonResp)
	w.WriteHeader(http.StatusOK)
}

// func ExtractTraceInfoMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
