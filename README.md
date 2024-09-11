# Simple-REST Web Framework

Simple-REST is a web framework written in [Go](https://go.dev/). It is a simple wrapper around the [net/http](https://pkg.go.dev/net/http) standard library. 
It has been developped to be as simple as possible while providing the most usefull utilities.

## Getting started

### Prerequisites

Simple-REST requires [Go](https://go.dev/) version [1.22](https://go.dev/doc/devel/release#go1.22.0) or above.

### Getting Simple-REST

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```sh
import "github.com/bragdond/simple-rest"
```

Alternatively, use `go get`:

```sh
go get -u github.com/bragdond/simple-rest
```

### Running Simple-REST

A basic example:

```go
package main

import (
  "net/http"

  "github.com/simple-rest"
)

func handleHello(w http.ResponseWriter, r *http.Request, p simplerest.Parameters) error {
    w.WriteHeader(http.StatusAccepted)
    _, err :=w.Write([]byte(data))
    return err
}

func main() {
    server := simplerest.NewServer("localhost", 8080)
    server.HandleFunc("/hello", nil, handleHello, http.MethodGet)
	go server.Serve()
}
```

To run the code, use the `go run` command, like:

```sh
$ go run cmd/main.go
```

Then visit [`localhost:8080/hello`](http://localhost:8080/hello) in your browser to see the response!


### Contributing

I welcome contributions! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Acknowledgements

- [Go](https://go.dev/)
- [net/http](https://pkg.go.dev/net/http)

### Contact

For any questions or suggestions, feel free to open an issue.
