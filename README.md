# crawler

A simple web crawler to explore concurrency in Go.

## Running

Flags:

```
  -limit int
        max number of pages to visit (default 10)
  -start_url string
        starting point
  -timeout int
        timeout in seconds (default 5)
  -workers int
        number of concurrent workers (default 1)
```

Example:

```bash
$ dep ensure
$ go run cmd/crawler/main.go -start_url https://example.com -workers 5
```

## Testing

```bash
$ go test -race ./...
```

## TODO

* Switch to `vgo`
* More extensive testing
