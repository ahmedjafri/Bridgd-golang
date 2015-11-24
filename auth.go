package main

import (
    "encoding/base64"
    "net/http"
    "strings"
)

type handler func(w http.ResponseWriter, r *http.Request)

// this is using basic auth which is very bad
// TODO (ajafri): use a more meaningful authentication scheme here
func BasicAuth(pass handler) handler {

    return func(w http.ResponseWriter, r *http.Request) {
        if _, ok := r.Header["Authorization"]; !ok {
            http.Error(w, "Authorization is required", http.StatusUnauthorized)
            return
        }

        auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

        if len(auth) != 2 || auth[0] != "Basic" {
            http.Error(w, "bad syntax", http.StatusBadRequest)
            return
        }

        payload, _ := base64.StdEncoding.DecodeString(auth[1])
        pair := strings.SplitN(string(payload), ":", 2)

        if len(pair) != 2 || !Validate(pair[0], pair[1]) {
            http.Error(w, "authorization failed", http.StatusUnauthorized)
            return
        }

        pass(w, r)
    }
}

func Validate(username, password string) bool {
    if username == "username" && password == "password" {
        return true
    }
    return false
}