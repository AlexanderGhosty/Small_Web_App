package handlers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"

    "small_web_app/internal/models"
)

var jwtKey = []byte("SUPER_SECRET_KEY")

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var creds struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid data", http.StatusBadRequest)
        return
    }

    var user models.User
    var hashedPassword string
    err := db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", creds.Email).
        Scan(&user.ID, &user.Name, &user.Email, &hashedPassword)
    if err != nil {
        http.Error(w, "Incorrect login or password", http.StatusUnauthorized)
        return
    }

    // Compare the hash
    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password)); err != nil {
        http.Error(w, "Incorrect login or password", http.StatusUnauthorized)
        return
    }

    // Generate token
    expirationTime := time.Now().Add(60 * time.Minute)
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// Example: can be used if token needs to be manually parsed
func ParseToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    return claims, nil
}
