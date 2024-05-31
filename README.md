### go-server

I've learned that go [1.22 released with "Enhanced routing patterns"](https://tip.golang.org/doc/go1.22), allowing you
to set up a fairly complex server with different routes and methods and stuff with just stdlib:

```go
package main

import (
	"fmt"
	"html"
	"net/http"
)

func HandleEndpoint(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
func main() {
	http.HandleFunc("GET /{$}", HandleEndpoint)
	http.HandleFunc("POST /api/todos/{id}", HandleEndpoint)
	http.HandleFunc("GET /todos", HandleEndpoint)
	http.HandleFunc("PUT /api/todos", HandleEndpoint)

	_ = http.ListenAndServe("8080", nil) 
}
```

Kinda exiting, huh? So I decided to try it.

So this is a project of me "trying it".
It has no other dependencies or libs, other than `htmx` for less javascript (which seems to be all the rage these days).
Right now it features:

- no dependencies
- routing
- middleware
- persistent database in its simplest form - a json file
- CRUD
- 404 page
- redirect to 404 on unknown urls
- htmx partials endpoint protection (redirect to 404 if not htmx request)
- config.json
- Two entities with common ID trait
- "Pure DI" (not my name: https://blog.ploeh.dk/2014/06/10/pure-di/)
- basic auth
- Docker
- Watch build size with `make`
