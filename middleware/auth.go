package middleware

import (
	"context"
	"net/http"
	"sghss-backend/utils"
	"strings"
)

// AuthMiddleware verifica se o usuário está autenticado
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obter o token do header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Token de autorização necessário"}`, http.StatusUnauthorized)
			return
		}

		// Verificar se o header tem o formato "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, `{"error": "Formato de token inválido"}`, http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Validar o token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, `{"error": "Token inválido"}`, http.StatusUnauthorized)
			return
		}

		// Adicionar as informações do usuário ao contexto da requisição
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_perfil", claims.Perfil)

		// Continuar para o próximo handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware verifica se o usuário é administrador
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		perfil := r.Context().Value("user_perfil")
		if perfil != "admin" {
			http.Error(w, `{"error": "Acesso negado. Apenas administradores"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

