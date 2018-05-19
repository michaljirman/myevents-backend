package rest

import (
	"net/http"

	"github.com/michaljirman/myevents-backend/src/lib/msgqueue"
	"github.com/michaljirman/myevents-backend/src/lib/persistence"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ServeAPI(endpoint, tlsendpoint string, dbhandler persistence.DatabaseHandler, eventEmmiter msgqueue.EventEmitter) (chan error, chan error) {
	handler := newEventHandler(dbhandler, eventEmmiter)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsrouter.Methods("GET").Path("/{eventID}").HandlerFunc(handler.oneEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	locationRouter := r.PathPrefix("/locations").Subrouter()
	locationRouter.Methods("GET").Path("").HandlerFunc(handler.allLocationsHandler)
	locationRouter.Methods("POST").Path("").HandlerFunc(handler.newLocationHandler)

	httpErrChan := make(chan error)
	httptlsErrChan := make(chan error)

	server := handlers.CORS()(r)
	go func() {
		httpErrChan <- http.ListenAndServeTLS(tlsendpoint, "/root/certs/cert.pem", "/etc/ssl/private/key.pem", server)
	}()
	go func() {
		httpErrChan <- http.ListenAndServe(endpoint, server)
	}()
	return httpErrChan, httptlsErrChan
}
