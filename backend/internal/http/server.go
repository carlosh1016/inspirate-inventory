package http

import (
	"fmt"
	"net/http"
	"time"
)

// NewServer builds an *http.Server with sane timeouts for handler, listening
// on port.
func NewServer(port int, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
