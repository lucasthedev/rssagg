package main

import (
	"fmt"
	"net/http"

	"github.com/lucasthedev/rssagg/internal/auth"
	"github.com/lucasthedev/rssagg/internal/database"
)

type autheHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middleware(handler autheHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)

		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUsersByApiKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Erro ao buscar usuário com api key no banco de dados: %v", err))
			return
		}
		handler(w, r, user)
	}
}
