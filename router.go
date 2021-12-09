package goapp

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rodriez/restface"
	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type myRouter struct {
	*mux.Router
}

func (r *myRouter) SecureFunc(pattern string, fn func(http.ResponseWriter, *http.Request)) *mux.Route {
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		defer HandlePanic(w)

		if err := Authenticator.CheckJWT(w, r); err == nil {
			fn(w, r)
		}
	}

	return r.Router.HandleFunc(pattern, wrapper)
}

func (r *myRouter) SNSFunc(pattern string, fn func(http.ResponseWriter, *http.Request)) *mux.Route {
	subscriptionConfirm := func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]string

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			logrus.Error(err)
		}

		if response, err := http.Get(payload["SubscribeURL"]); err != nil {
			logrus.Error(err)
		} else {
			presenter := restface.Presenter{Writer: w}
			presenter.Present(http.StatusOK, response)
		}
	}

	wrapper := func(w http.ResponseWriter, r *http.Request) {
		defer HandlePanic(w)

		if r.Header.Get("x-amz-sns-message-type") == "SubscriptionConfirmation" {
			subscriptionConfirm(w, r)
			return
		}

		fn(w, r)
	}

	return r.Router.HandleFunc(pattern, wrapper)
}

var Router *myRouter

func InitRouter() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/ping", Ping).Methods(http.MethodGet)
	rtr.HandleFunc("/status/ping", Ping).Methods(http.MethodGet)
	mux.CORSMethodMiddleware(rtr)

	Router = &myRouter{rtr}
	Router.Use(otelmux.Middleware(os.Getenv("TRACE_ID")))
}
