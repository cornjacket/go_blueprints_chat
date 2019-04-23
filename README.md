# go_blueprints_chat
go programming blue prints chat app

p14 use of lazy loading via the sync.Once mechanism. Prior to serving the first request for the template, the template will be compiled. useful in cases where there is a lot of preliminary calculations that are not necessary on program start.

p.18 use of gorilla/websocket package to implement websockets

p.18-19 room module's use of join and leave channels allow safe sharing of common map for concurrent goroutines.

p.20 shows a different way to use select. the forward channel receives a message which causes the room to iterate through all the clients and send the message to each. however the code uses a select statement with a case that sends the msg to the client's send channel. the default case handles removing the client from the clients map. the defeault case behaves like an exception case. however it is unclear why a select statement is being used since the purpose is not to wait for a msg to be received but instead to send a msg to the client. why not simply have a if (send) command with the error condition handled via an else. This looks like a special use of the select statement.

p.21 add ServeHTTP handler for room type. room can now act as a handler. code upgrades the HTTP connection using the websocket.Upgrader type.

p22 update main to reflect /room path as well as kicking room off

Tested on localhost:8082 and gcloud:8082.

p.38 add trace package using interfaces. allows tracing during test that is disabled during deployment by having defeault nil tracer.

p.39 take away - using web sockets allows communication with the client without messy polling. Research this.

p.48 /login page that gives various authentication paths. Dynamic path functionality that is stubbed, /auth/login/provider.

p.54 w.Header.Set looks like a typo. It should be w.Header().Set

p.58 added authentication using google. facebook and github are non-functional paths.

p.64 add username and when into message forwarded to each client. Defined chan of message type to forward from room to client.
