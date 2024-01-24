package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"service/gateway/dto"
	g "service/gateway/global"

	"service/pkg/errors"
	"service/pkg/translator"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		translate := ctx.Value(g.TranslateKey).(translator.TranslatorFunc)
		defer func() {
			errInterface := recover()
			if errInterface == nil {
				return
			}
			if err, ok := errInterface.(error); ok && errors.IsServerError(err) {
				code, action, message, _, errors := errors.HttpError(err)
				res := dto.PanicResponse{
					Message: translate(message),
					Code:    code,
					Action:  action,
					Errors:  errors,
				}
				resBytes, _ := json.Marshal(res)
				if g.CFG.Debug {
					log.Println(err)
				}
				if r.Header.Get("timeout") == "yes" {
					return
				}
				w.WriteHeader(code)
				w.Write(resBytes)
				if code == http.StatusRequestTimeout {
					r.Header.Set("timeout", "yes")
				}
			} else {
				stack := string(debug.Stack())
				g.Logger.Panic(errInterface, r, stack)
				res := dto.PanicResponse{
					Message: translate("InternalServerError"),
					Action:  int(errors.Report),
					Code:    http.StatusInternalServerError,
					Errors:  nil,
				}
				resBytes, _ := json.Marshal(res)
				w.WriteHeader(res.Code)
				w.Write(resBytes)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
