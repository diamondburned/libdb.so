// Package backend implements the backend for libdb.so.
package backend

import "net/http"

// Params defines the parameters for the backend.
type Params struct {
	health.HealthService
}

func New() http.Handler {}
