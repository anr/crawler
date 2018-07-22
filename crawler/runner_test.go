package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var pages = map[string]string{
	"/a": `<!DOCTYPE html>
		<html>
			<head><title>title a</title></head>
			<body>
				<a href="/b">b</a>
				<a href="/c">c</a>
			</body>
		</html>`,
	"/b": `<!DOCTYPE html>
		<html>
			<head><title>title b</title></head>
			<body>
				<a href="/a">a</a>
				<a href="http://127.1.1.1">external</a>
			</body>
		</html>`,
	"/c": `<!DOCTYPE html>
		<html>
			<head></head>
			<body>
				<a href="/missing">missing</a>
			</body>
		</html>`,
}

func mustGetAbsURL(server *httptest.Server, myURL string) string {
	u, _, err := absURL(server.URL, myURL)
	if err != nil {
		panic(err)
	}

	return u
}

func wantSiteMap(server *httptest.Server) map[string]string {
	return map[string]string{
		mustGetAbsURL(server, "/a"): "title a",
		mustGetAbsURL(server, "/b"): "title b",
		mustGetAbsURL(server, "/c"): "",
	}
}

// makeHandler creates a http.HandlerFunc to serve the pages passed.
// pages map paths to html content.
func makeHandler(pages map[string]string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received request for page %q", r.URL)
			page, ok := pages[r.URL.EscapedPath()]
			if !ok {
				log.Println("Page was not found")
				http.Error(w, "page not found", http.StatusNotFound)
				return
			}

			fmt.Fprint(w, page)
		})
}

func TestEndToEnd(t *testing.T) {
	server := httptest.NewServer(makeHandler(pages))
	defer server.Close()

	startURL := mustGetAbsURL(server, "/a")

	for _, workers := range []int{1, 10} {
		crawler := &Crawler{
			StartURL: startURL,
			Workers:  workers,
			Client:   server.Client(),
			Logger:   log.New(os.Stdout, "", log.LstdFlags),
		}

		want := wantSiteMap(server)

		if got := crawler.Run(); !reflect.DeepEqual(want, got) {
			t.Errorf("want %s, got %s", want, got)
		}
	}
}
