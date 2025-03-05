package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("lolit@456a'")


func GenerateJWT(email, role string) (string, error) {
    
    claims := &jwt.MapClaims{
        "email": email,
        "role":  role,
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
    }


    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(JwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
