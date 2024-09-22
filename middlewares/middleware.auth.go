package middlewares

import (
	"net/http"
	"os"

	"github.com/restuwahyu13/discovery-api/helpers"
	"github.com/restuwahyu13/discovery-api/packages"
)

func AuthToken(next http.Handler) http.Handler {
	var (
		res        *helpers.Response   = new(helpers.Response)
		asymmetric *helpers.Asymmetric = new(helpers.Asymmetric)
		rsa        helpers.IRsaCrypto  = helpers.NewRsa()
		password   string              = os.Getenv("INTERNAL_PRIVATE_KEY_PASSWORD")
	)

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		res.StatCode = http.StatusUnauthorized

		asymmetric.ClientID = r.Header.Get("X-CLIENT-ID")
		asymmetric.ClientKey = r.Header.Get("X-CLIENT-KEY")
		asymmetric.PrivateKey = r.Header.Get("X-PRIVATE-KEY")

		if asymmetric.ClientID == "" || asymmetric.ClientKey == "" || asymmetric.PrivateKey == "" {
			res.ErrCode = "UNAUTHORIZED"
			res.ErrMsg = "Invalid credentials"

			helpers.ApiResponse(rw, res)
			return
		}

		if err := rsa.Asymmetric(asymmetric, password); err != nil {
			res.ErrCode = "UNAUTHORIZED"
			res.ErrMsg = "Invalid credentials"

			defer packages.Logrus("error", err)
			helpers.ApiResponse(rw, res)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
