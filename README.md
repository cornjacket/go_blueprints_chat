# go_blueprints_chat
go programming blue prints chat app

p14 use of lazy loading via the sync.Once mechanism. Prior to serving the first request for the template, the template will be compiled. useful in cases where there is a lot of preliminary calculations that are not necessary on program start.

p.18 use of gorilla/websocket package to implement websockets
