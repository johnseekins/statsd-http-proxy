package routehandler

import (
	"fmt"
	"net/http"
)

// Handle heartbeat request
func (routeHandler *RouteHandler) HandleHeartbeatRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
