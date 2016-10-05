package middlewares

import (
	"fmt"
	"github.com/newrelic/go-agent"
	"net/http"
	"strings"
)

// Nragent is a middleware wich send information on Newrelic for all requests
type Nragent struct {
	Application *newrelic.Application
	Transaction *newrelic.Transaction
}

// NewNragent returns an initialized NrAgent.
func NewNragent(appname string, secretkey string) *Nragent {
	config := newrelic.NewConfig(appname, secretkey)
	config.Enabled = true
	app, _ := newrelic.NewApplication(config)
	return &Nragent{Application: &app}
}

func (n *Nragent) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	path := r.URL.Path
	apiPath := strings.Split(path, "/")[1]
	host := r.Host
	txn := ((*n.Application).StartTransaction(host+"/"+apiPath, rw, r)).(newrelic.Transaction)
	extSeg := external(txn, r)
	n.Transaction = &txn
	defer (*n.Transaction).End()
	defer changeName(*extSeg, r)
	next(rw, r)
}

func changeName(tr newrelic.ExternalSegment, req *http.Request) {
	if reqbakNameHdr := req.Header["X-Traefik-backName"]; len(reqbakNameHdr) == 1 {
		backendServer := reqbakNameHdr[0]
		//fmt.Println(backendServer + req.RequestURI)
		tr.URL = backendServer + req.RequestURI
	}
	tr.End()
}

func external(txn newrelic.Transaction, req *http.Request) *newrelic.ExternalSegment {
	extSeg := newrelic.ExternalSegment{
		StartTime: newrelic.StartSegmentNow(txn),
		//Request:   req,
		URL: "",
	}
	return &extSeg
}
