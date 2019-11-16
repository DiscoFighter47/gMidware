package gmidware

import (
	"errors"
	"log"
	"net/http"

	gero "github.com/DiscoFighter47/gEro"
	gson "github.com/DiscoFighter47/gSON"
)

// Recoverer ...
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Error Occurred:", err)
				switch err := err.(type) {
				case *gero.APIerror:
					gson.ServeError(w, err)
				case error:
					gson.ServeError(w, gero.NewAPIerror("Internal Server Error", http.StatusInternalServerError, err))
				case string:
					gson.ServeError(w, gero.NewAPIerror("Internal Server Error", http.StatusInternalServerError, errors.New(err)))
				default:
					panic(err)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
