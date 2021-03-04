# middleware
A handy tool for chaining http.Handler middlewares in Golang

```
import "github.com/hashibuto/middleware"
```

## The problem

Middleware chaining involves one middleware calling another, which in turn calls another, and so on.  This is both irritating to write, as well as annoying to reuse.

## The solution

A very simple module which auto-wraps middlewares so that they can be conveniently called, and
wrap in the order in which they occur.

## Example

Consider these two middlewares:

```
func middleA(next http.Handler) http.Handler {
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      // Do something before the request
      ...
      next.ServeHTTP(w, r)
      ...
      // Do something after the request
    },
  )
}

func middleB(next http.Handler) http.Handler {
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      // Do something before the request
      ...
      next.ServeHTTP(w, r)
      ...
      // Do something after the request
    },
  )
}
```

If we want to have these two middlewares handle requests, we'd have to do the following:

```
  mux := http.NewServeMux()
  mux.Handle("/", http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      // Handle the root request
    },
  ))

  http.ListenAndServe(
    "0.0.0.0:8084",
    middleA(
      middleB(
        mux
      )
    )
  )
```

Notice how each handler wraps the next.

With the middleware module, this becomes easier to read and more convenient:

```
  mux := http.NewServeMux()
  mux.Handle("/", http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      // Handle the root request
    },
  ))

  http.ListenAndServe(
    "0.0.0.0:8084",
    middleware.Cascade(
      mux,
      middleA,
      middleB,
    )
  )
```

Similarly, middlewares can be stored and reused.  Consider this example:

```
  mux := http.NewServeMux()

  myMiddlewares = []MiddleWare{
    middleA,
    middleB,
  }

  rootHandler := http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      // Handle the root request
    },
  )

  helloHandler := http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      // Handle the hello request
    },
  )

  mux.Handle("/", middlewares.Cascade(rootHandler, myMiddlewares...))
  mux.Handle("/hello", middlewares.Cascade(helloHandler, myMiddlewares...))

  http.ListenAndServe(
    "0.0.0.0:8084",
    mux,
  )
```
