package middleware

import (
    "context"
    "net/http"
    "strings"

    "small_web_app/internal/handlers"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Отсутствует токен авторизации", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Неверный формат заголовка Authorization", http.StatusUnauthorized)
            return
        }
        tokenStr := parts[1]

        claims, err := handlers.ParseToken(tokenStr)
        if err != nil {
            http.Error(w, "Невалидный или просроченный токен", http.StatusUnauthorized)
            return
        }

        // Пробрасываем userID в контекст запроса
        ctx := context.WithValue(r.Context(), "userID", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
