package goapp

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/rodriez/restface"
)

var Authenticator *jwtmiddleware.JWTMiddleware

var JWTExtraValidation func(token *jwt.Token) error

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func InitAuth() {
	Authenticator = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			aud := os.Getenv("JWT_AUD")
			if checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false); !checkAud {
				return token, errors.New("invalid audience")
			}

			iss := os.Getenv("JWT_ISS")
			if checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false); !checkIss {
				return token, errors.New("invalid issuer")
			}

			if JWTExtraValidation != nil {
				if err := JWTExtraValidation(token); err != nil {
					return token, err
				}
			}

			if cert, err := getPemCert(token); err == nil {
				result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
				return result, nil
			} else {
				return nil, err
			}
		},
		SigningMethod: jwt.SigningMethodRS256,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			presenter := restface.Presenter{Writer: w}
			presenter.PresentError(restface.Unauthorized(err))
		},
	})
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	jwks, err := getJwks()
	if err != nil {
		return "", err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func getJwks() (*Jwks, error) {
	resp, err := http.Get(os.Getenv("JWKS_URL"))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var jwks = Jwks{}
	if err = json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}
