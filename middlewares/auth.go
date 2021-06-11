package middlewares

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/kluddizz/maintenance-rest-service/models"
)

func AuthMiddleWare(next httprouter.Handle) httprouter.Handle {
  return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
    res := models.NewJsonResponse(w)
    authHeader := req.Header.Get("authorization")

    if authHeader != "" {
      bearerToken := strings.Split(authHeader, " ")

      if len(bearerToken) == 2 {
        token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
          if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
          }

          privateKey, err := ioutil.ReadFile("./private.key")

          if err != nil {
            return nil, fmt.Errorf("Error while reading private key")
          }

          return []byte(privateKey), nil
        })

        if err == nil && token.Valid {
          ctx := context.WithValue(req.Context(), "auth", token.Claims)
          next(w, req.WithContext(ctx), p)
        } else {
          res.Code = 400
          res.Content = "Invalid authorization token"
          res.Send()
        }
      }
    } else {
      res.Code = 400
      res.Content = "An authorization header is required"
      res.Send()
    }
  }
}
