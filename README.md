### go-server

I've learned that go [1.22 released with "Enhanced routing patterns"](https://tip.golang.org/doc/go1.22), allowing you
to set up a fairly complex server with different routes and methods and stuff with just stdlib:

```go
package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("GET /{$}", IndexGetHandler)
	http.HandleFunc("POST /api/todos/{id}", ApiTodosPostHandler)
	http.HandleFunc("GET /todos", TodosGetHandler)
	http.HandleFunc("PUT /api/todos", ApiTodosPutHandler)

	_ = http.ListenAndServe("8080", nil)
}
```

Kinda exiting, huh? So I decided to try it.

So this is a project of me "trying it".
It has no other dependencies or libs, other than `htmx` for less javascript (which seems to be all the rage these days).
Right now it features:

- Routing
- Middleware
- Persistent database in its simplest form - a json file
- CRU(d)
