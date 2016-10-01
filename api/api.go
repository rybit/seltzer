package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/guregu/kami"
	"github.com/pborman/uuid"
	"github.com/rs/cors"
	"golang.org/x/net/context"

	"github.com/rybit/config_example/conf"
)

type API struct {
	log     *logrus.Entry
	config  *conf.Config
	handler http.Handler
}

type JWTClaims struct {
	jwt.StandardClaims
	Name    string
	IsAdmin bool
}

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

func NewAPI(config *conf.Config) *API {
	api := &API{
		log:    logrus.WithField("component", "api"),
		config: config,
	}

	k := kami.New()
	k.Use("/", api.populateConfig)
	k.Get("/", hello)

	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	api.handler = corsHandler.Handler(k)
	return api
}

func (a API) populateConfig(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	reqID := uuid.NewRandom().String()
	log := a.log.WithFields(logrus.Fields{
		"request_id": reqID,
		"method":     r.Method,
		"path":       r.URL.Path,
	})

	ctx = setRequestID(ctx, reqID)
	ctx = setStartTime(ctx, time.Now())
	ctx = setConfig(ctx, a.config)

	token, err := extractToken(a.config.JWTSecret, r)
	if err != nil {
		log.WithError(err).Info("Failed to parse token")
		return nil
	}

	if token == nil {
		log.Debug("Making unauthenticated request")
	} else {
		claims := token.Claims.(*JWTClaims)
		log = log.WithField("is_admin", claims.IsAdmin)
		ctx = setAdminFlag(ctx, claims.IsAdmin)
	}

	ctx = setLogger(ctx, log)
	return ctx
}

func extractToken(secret string, r *http.Request) (*jwt.Token, *HTTPError) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, nil
	}

	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return nil, httpError(http.StatusBadRequest, "Bad authentication header")
	}

	token, err := jwt.ParseWithClaims(matches[1], &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, httpError(http.StatusUnauthorized, "Invalid Token")
	}

	claims := token.Claims.(*JWTClaims)
	if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
		return nil, httpError(http.StatusUnauthorized, fmt.Sprintf("Token expired at %v", time.Unix(claims.StandardClaims.ExpiresAt, 0)))
	}
	return token, nil
}

func hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	version := getVersion(ctx)
	sendJSON(w, http.StatusOK, map[string]string{
		"version":     version,
		"application": "seltzer",
	})
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.Encode(obj)
}
