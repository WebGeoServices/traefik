package middlewares

import (
	"github.com/newrelic/go-agent"
	"net/http"
)

// Nragent is a middleware wich send information on Newrelic for all requests
type Nragent struct {
	Application *newrelic.Application
	Transaction *newrelic.Transaction
}

// NewNrAgent returns an initialized NrAgent.
func NewNragent(appname string, secretkey string) *Nragent {
	config := newrelic.NewConfig(appname, secretkey)
	config.Enabled = true
	app, _ := newrelic.NewApplication(config)
	return &Nragent{Application: &app}
}

func (n *Nragent) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	txn := ((*n.Application).StartTransaction(r.URL.Path, rw, r)).(newrelic.Transaction)
	n.Transaction = &txn
	defer (*n.Transaction).End()
	next(rw, r)
}
