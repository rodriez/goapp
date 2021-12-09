package goapp

import (
	"fmt"
	"net/http"

	"github.com/rodriez/restface"
)

func HandlePanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		err := fmt.Errorf("%v", r)

		presenter := restface.NewPresenter(w)
		presenter.PresentError(restface.InternalError(err.Error()))
	}
}
