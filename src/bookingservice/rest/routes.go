package rest

import (
	"net/http"

	"github.com/gorilla/handlers"

	"github.com/michaljirman/myevents-backend/src/lib/msgqueue"
	"github.com/michaljirman/myevents-backend/src/lib/persistence"

	"github.com/gorilla/mux"
)

func ServeAPI(endpoint, tlsendpoint string, database persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) (chan error, chan error) {
	r := mux.NewRouter()
	r.Methods("POST").Path("/events/{eventID}/bookings").Handler(&CreateBookingHandler{eventEmitter, database})

	httpErrChan := make(chan error)
	httptlsErrChan := make(chan error)

	server := handlers.CORS()(r)
	// go func() {
	// 	httpErrChan <- http.ListenAndServeTLS(tlsendpoint, "/root/certs/cert.pem", "/etc/ssl/private/key.pem", server)
	// }()
	go func() {
		httpErrChan <- http.ListenAndServe(endpoint, server)
	}()
	return httpErrChan, httptlsErrChan
}
