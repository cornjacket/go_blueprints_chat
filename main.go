package main

import (
  "flag"
  "log"
  "net/http"
  "text/template"
  "path/filepath"
  "sync"
  //"os"
  //"github.com/cornjacket/trace"
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/signature"
  "github.com/stretchr/objx"
  "github.com/stretchr/gomniauth/providers/google"
)

// templ represents a signle template
type templateHandler struct {
  once     sync.Once
  filename string
  templ    *template.Template
}

// ServeHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.once.Do(func() {
    t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
  })
  data := map[string]interface{}{
    "Host": r.Host,
  }
  if authCookie, err := r.Cookie("auth"); err == nil {
    data["UserData"] = objx.MustFromBase64(authCookie.Value)
  }
  t.templ.Execute(w, data)
}

func main() {
  var addr = flag.String("addr", ":8082", "The addr of the application.")
  flag.Parse() // parse the flags
  // setup gomniauth
  gomniauth.SetSecurityKey(signature.RandomKey(64))
  gomniauth.WithProviders( google.New("1039434619377-fnfcgtdrhj82ssto8q95s9jqhfms5d73.apps.googleusercontent.com", "hm12hOSAM8HeAElysMiGx-vd", "http://localhost:8082/auth/callback/google") )

  r := newRoom(UseFileSystemAvatar)
  //r := newRoom(UseGravatar)
  //r := newRoom(UseAuthAvatar)
  //r.tracer = trace.New(os.Stdout) // only used for test tracing
  http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
  http.Handle("/login", &templateHandler{filename: "login.html"})
  http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
      Name: "auth",
      Value: "",
      Path: "/",
      MaxAge: -1,
    })
    w.Header().Set("Location", "/chat")
    w.WriteHeader(http.StatusTemporaryRedirect)
  })
  http.HandleFunc("/auth/", loginHandler)
  http.Handle("/room", r)
  http.Handle("/upload", &templateHandler{filename: "upload.html"})
  http.HandleFunc("/uploader", uploaderHandler)
  http.Handle("/avatars/",
    http.StripPrefix("/avatars/",
      http.FileServer(http.Dir("./avatars"))))
  // get the room going
  go r.run()
  // start the web server
  log.Println("Starting web server on", *addr)
  if err := http.ListenAndServe(*addr, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
