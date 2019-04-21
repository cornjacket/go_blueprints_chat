package main

import (
  "flag"
  "log"
  "net/http"
  "text/template"
  "path/filepath"
  "sync"
  "os"
  "github.com/cornjacket/trace"
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
  t.templ.Execute(w, r)
}

func main() {
  var addr = flag.String("addr", ":8082", "The addr of the application.")
  flag.Parse() // parse the flags
  r := newRoom()
  r.tracer = trace.New(os.Stdout)
  http.Handle("/", &templateHandler{filename: "chat.html"})
  http.Handle("/room", r)
  // get the room going
  go r.run()
  // start the web server
  log.Println("Starting web server on", *addr)
  if err := http.ListenAndServe(*addr, nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
