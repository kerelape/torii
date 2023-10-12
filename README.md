# TORII

Superstructure on top of `http.Handler` the way it should be.

## Example

### `http`

```go
package main

import "net/http"

func main() {
    var service interface{
        SomeData(context.Context) ([]byte, error)
    }
    http.ListenAndServe(
        ":8080",
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodGet {
                http.Error(
                    w,
                    http.StatusText(http.StatusMethodNotAllowed),
                    http.StatusMethodNotAllowed,
                )
                return
            }
            data, dataError := service.SomeData(r.Context())
            if dataError != nil {
                http.Error(
                    w,
                    http.StatusText(http.StatusInternalServerError),
                    http.StatusInternalServerError,
                )
                return
            }
            w.WriteHeader(http.StatusOK)
            if _, err := w.Write(data); err != nil {
                log.Println(err)
            }
        }),
    )
}

```

### `torii`

```go
package main

import (
    "net/http"

    "github.com/kerelape/torii"
)

func main() {
    var service interface{
        SomeData(context.Context) ([]byte, error)
    }
    http.ListenAndServe(
        ":8080",
        torii.Handler(
            torii.Func(func(ctx context.Context, request torii.Request) torii.Response {
                if request.Method != torii.MethodGet {
                    return torii.Response{
                        Status: torii.StatusMethodNotAllowed,
                    }
                }
                data, dataError := service.SomeData(ctx)
                if dataError != nil {
                    return torii.Response{
                        Status: torii.StatusInternalServerError,
                    }
                }
                return torii.Response{
                    Status: torii.StatusOK,
                    Body:   io.NopCloser(bytes.NewReader(data)),
                }
            },
            torii.HandlerOptionWithLogger(log.Default()),
        ),
    )
}

```
