# Gouter: Yet another router for Go


### Why another router?

These are my reasons to write this:
 1. Learning Go
 2. Create simpler router
 3. Making a router is fun and easy
 4. To fight with my burnout (-_-)
 
### Todo

- [x] Middlewares
- [x] Support params in URL 
- [x] Send custom type to the handlers (Context)
- [ ] More tests
- [ ] Benchmark
- [ ] Documentation
- [ ] Get URL by route name
- [ ] attach multiple context to the Request

Installation

`go get github.com/smoqadam/gouter`


## How to use Gouter

### Simple 

```go

func main() {
  r := router.NewRouter()

  r.GET("/", indexHandler).Name("index")

  http.ListenAndServe(":3000", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
   fmt.FprintF(w, "hello, world!")
}

```

### Route with Params

You can have dynamic path and send the paramters to the handler:

```go

func main() {
  r := router.NewRouter()

  r.GET("/user/{user}", userHandler).
      Name("index").
      Where("user", "[a-z0-9]+")

  http.ListenAndServe(":3000", r)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
   vars := router.Vars(r)
   fmt.Fprintf(w, "hello, world!", vars["user"])
}

```

####  Middlewares

Middleware method receives one or more middleware function and execute them before the final handler. The middleware function receives a `http.handler` and return a `http.handler`. By calling `next.ServeHTTP(w, r)` at the end of your handler it will run the next middleware or the final handler.

```go

func main() {
  r := router.NewRouter()

  r.GET("/user/{user}", userHandler).
      Name("index").
      Where("user", "[a-z0-9]+").
      Middleware(mid1, mid2)

  http.ListenAndServe(":3000", r)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
   vars := router.Vars(r)
   fmt.FprintF(w, "hello, world!", vars["user"])
}

func mid1(next http.Handler) http.Handler {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
      fmt.Println("from middleware 1")
      next.ServeHTTP(w, r) // call another middleware or the final handler
  });
}

func mid2(next http.Handler) http.Handler {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
      fmt.Println("from middleware 2")
      w.WriteHeader(http.StatusBadRequest) // send 400 status code to the client
  });
}
```


### Attach a context to the Request

Sometimes we need to send an struct to our handler such as controller or model obejct. `With` method in router instance recieves an `interface` and send it through `http.Request` to the handler.


```go
package main

import (
	"fmt"
	"net/http"

	router "github.com/smoqadam/gouter"
)

type App struct {
	Name string
}

func main() {
	r := router.NewRouter()
	app := &App{Name: "Gouter"}
	r.With(app)
	r.GET("/user/{user}", userHandler).
		Name("index").
		Where("user", "[a-z0-9]+")

	http.ListenAndServe(":3000", r)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// We can access to *App like this
	app := router.Context(r).(*App)
	fmt.Fprintf(w, "hello, %s!", app.Name)
}
```

