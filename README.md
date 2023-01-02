# Path Proxy

Super simple way to define a proxy based on matching a url path.

## Example

In this example we:

1. start a proxy server on port `:3333`
2. send all traffic that matches `^/api` (starts with /api) to `http://localhost:7788`
3. send everything else (`/`) to `http://localhost:3000`

```go
package main

import (
    "fmt"
	"net/http"

    "github.com/emehrkay/path_proxy"
)

func main() {
    var err error
	handler := &PathProxyHandler{}
	err = handler.ProxyDefintions(ProxyDefSet{
		{
			"/",
			"http://localhost:3000",
		},
		{
			"^/api",
			"http://localhost:7788",
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("server running on port 3333")
	http.ListenAndServe(":3333", handler)
}
```

> Route definitions are automatically ordered from longest to shortest.
> This is also a simple regular expression handler `handler.Handle("/some/regex$", anActualHandler)`
