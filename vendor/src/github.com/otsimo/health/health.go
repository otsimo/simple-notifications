package health

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// This file copied from github.com/coreos/pkg

const (
	JSONContentType = "application/json"
)

// Checkables should return nil when the thing they are checking is healthy, and an error otherwise.
type Checkable interface {
	Healthy() error
}

func New(checks ...Checkable) *Checker {
	return &Checker{
		Checks:           checks,
		UnhealthyHandler: DefaultUnhealthyHandler,
		HealthyHandler:   DefaultHealthyHandler,
	}
}

// Checker provides a way to make an endpoint which can be probed for system health.
type Checker struct {
	// Checks are the Checkables to be checked when probing.
	Checks []Checkable

	// Unhealthyhandler is called when one or more of the checks are unhealthy.
	// If not provided DefaultUnhealthyHandler is called.
	UnhealthyHandler UnhealthyHandler

	// HealthyHandler is called when all checks are healthy.
	// If not provided, DefaultHealthyHandler is called.
	HealthyHandler http.HandlerFunc
}

func (c Checker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	unhealthyHandler := c.UnhealthyHandler
	if unhealthyHandler == nil {
		unhealthyHandler = DefaultUnhealthyHandler
	}

	successHandler := c.HealthyHandler
	if successHandler == nil {
		successHandler = DefaultHealthyHandler
	}

	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := Check(c.Checks); err != nil {
		unhealthyHandler(w, r, err)
		return
	}

	successHandler(w, r)
}

func (c Checker) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	if err := Check(c.Checks); err != nil {
		log.Printf("health check failed: %v", err)
		return &healthpb.HealthCheckResponse{
			Status: healthpb.HealthCheckResponse_NOT_SERVING,
		}, nil
	}
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}

type UnhealthyHandler func(w http.ResponseWriter, r *http.Request, err error)

type StatusResponse struct {
	Status  string                 `json:"status"`
	Details *StatusResponseDetails `json:"details,omitempty"`
}

type StatusResponseDetails struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func Check(checks []Checkable) (err error) {
	errs := []error{}
	for _, c := range checks {
		if e := c.Healthy(); e != nil {
			errs = append(errs, e)
		}
	}

	switch len(errs) {
	case 0:
		err = nil
	case 1:
		err = errs[0]
	default:
		err = fmt.Errorf("multiple health check failure: %v", errs)
	}

	return
}

func DefaultHealthyHandler(w http.ResponseWriter, r *http.Request) {
	err := writeJSONResponse(w, http.StatusOK, StatusResponse{
		Status: "ok",
	})
	if err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}

func DefaultUnhealthyHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("health check failed: %v", err)
	writeErr := writeJSONResponse(w, http.StatusInternalServerError, StatusResponse{
		Status: "error",
		Details: &StatusResponseDetails{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	})
	if writeErr != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}

func writeJSONResponse(w http.ResponseWriter, code int, resp interface{}) error {
	enc, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", JSONContentType)
	w.WriteHeader(code)

	_, err = w.Write(enc)
	if err != nil {
		return err
	}
	return nil
}
