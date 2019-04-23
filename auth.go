package main

import (
  "net/http"
  "strings"
  "fmt"
  //"log"
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/objx"
)

type authHandler struct {
  next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
    // not authenticated
    w.Header().Set("Location", "/login")
    w.WriteHeader(http.StatusTemporaryRedirect)
  } else if err != nil {
    // some other error
    panic(err.Error())
  } else {
    // success - call the next handler
    h.next.ServeHTTP(w, r)
  }
}

func MustAuth(handler http.Handler) http.Handler {
  return &authHandler{next: handler}
}

// loginHandler handles the third-party process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
  segs := strings.Split(r.URL.Path, "/")
  if len(segs) < 4 {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "Auth action not supported. Insufficient params")
    return
  }
  action := segs[2]
  provider := segs[3]
  switch action {
  case "login":
    provider, err := gomniauth.Provider(provider)
    if err != nil {
      http.Error(w, fmt.Sprintf("Error when trying to get the provider %s: %s", provider, err), http.StatusBadRequest)
      return
    }
    loginUrl, err := provider.GetBeginAuthURL(nil, nil)
    if err != nil {
      http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s: %s", provider, err), http.StatusInternalServerError)
      return
    }
    w.Header().Set("Location", loginUrl)
    w.WriteHeader(http.StatusTemporaryRedirect)
  case "callback":
    provider, err := gomniauth.Provider(provider)
    if err != nil {
      http.Error(w, fmt.Sprintf("Error when trying to get the provider %s: %s", provider, err), http.StatusBadRequest)
      return
    }
    creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
    if err != nil {
      http.Error(w, fmt.Sprintf("Error when trying to complete auth for %s: %s", provider, err), http.StatusInternalServerError)
      return
    }
    user, err := provider.GetUser(creds)
    if err != nil {
      http.Error(w, fmt.Sprintf("Error when trying to get user from %s: %s", provider, err), http.StatusInternalServerError)
      return
    }
    fmt.Printf("auth: user.Name = %s\n", user.Name())
    authCookieValue := objx.New(map[string]interface{}{
      "name": user.Name(),
    }).MustBase64()
    http.SetCookie(w, &http.Cookie{
      Name: "auth",
      Value: authCookieValue,
      Path: "/"})
    w.Header().Set("Location", "/chat")
    w.WriteHeader(http.StatusTemporaryRedirect)
  default:
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "Auth action %s not supported", action)
  }
}
