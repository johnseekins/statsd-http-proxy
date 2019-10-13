package routehandler

import (
	"fmt"
	"net/http"
)

// Handle heartbeat request
func (routeHandler *routeHandler) handleHeartbeatRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
