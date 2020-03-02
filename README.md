# gouter
Yet another router for Go

### Why another router?

First of all I'm learning Go and IMO, learning by doing is the best way to learn. Secondly, I used gorilla/mux router and negroni 
middleware but I didn't like them. As a PHP/JS developer the gorolla/mux and negroni was not so stright forward for me. 
So I wrote this for two reason:

 1. Learning Go
 2. Create simpler router
 3. Making a router is fun and easy (-_-)

### Todo

- [x] Middlewares
- [x] Support params in URL 
- [x] Send custom type to the handlers (Context)
- [ ] More tests
- [ ] Documentation
- [ ] Get URL by route name
- [ ] 

Installation

`go get github.com/smoqadam/gouter`


## How to use Gouter

#### Simple 

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

#### Route with Params

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
   fmt.FprintF(w, "hello, world!", r["user"])
}

```

####  Middlewares

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
      next.ServeHTTP(w, r) // call another middleware or the final handler
  });
}
```


### Context
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
	app := router.Context(r).(*App)
	fmt.Fprintf(w, "hello, %s!", app.Name)
}
```




