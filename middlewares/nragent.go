package middlewares

import (
	"net/http"
	"github.com/newrelic/go-agent"
)


type Nragent struct {
	Application *newrelic.Application
	Transaction *newrelic.Transaction
}


func NewNragent(appname string, secretkey string) (*Nragent) {
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