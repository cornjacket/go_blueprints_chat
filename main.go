package main

import (
  "log"
  "net/http"
  "text/template"
  "path/filepath"
  "sync"
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
  t.templ.Execute(w, nil)
}

func main() {
  r := newRoom()
  http.Handle("/", &templateHandler{filename: "chat.html"})
  http.Handle("/room", r)
  // get the room going
  go r.run()
  // start the web server
  if err := http.ListenAndServe(":8082", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
