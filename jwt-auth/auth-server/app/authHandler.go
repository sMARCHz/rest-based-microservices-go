package app

import (
	"encoding/json"
	"net/http"

	"github.com/sMARCHz/rest-based-microservices-go-lib/logger"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/dto"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/services"
)

type AuthHandler struct {
	service services.AuthService
}

func (a AuthHandler) NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, "Handler not implemented...")
}

func (a AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logger.Error("Error while decoding login request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := a.service.Login(loginRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr)
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

func (a AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	// converting from Query to map type
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	appErr := a.service.Verify(urlParams)
	if appErr != nil {
		writeResponse(w, appErr.Code, unAuthorizedResponse(appErr.Message))
	} else {
		writeResponse(w, http.StatusOK, authorizedResponse())
	}
}

func (a AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshTokenRequest dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshTokenRequest); err != nil {
		logger.Error("Error while decoding refresh request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := a.service.Refresh(refreshTokenRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr)
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

func unAuthorizedResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      msg,
	}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{
		"isAuthorized": true,
	}
}
