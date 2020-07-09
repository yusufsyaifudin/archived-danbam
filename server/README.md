# server

Server just contains abstraction layer required by server to run.

So, every app inside `internal/app` which use REST API as interface, must implement
every abstraction in this `server` directory. This make us easy to call handler in `cmd`
when we need to build the binary.