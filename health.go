package goapp

import (
	"net/http"
	"time"

	"github.com/rodriez/restface"
)

func Ping(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
	}

	presenter := restface.Presenter{Writer: res}
	presenter.Present(http.StatusOK, body)
}
