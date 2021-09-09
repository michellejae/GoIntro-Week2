package main

// client side

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	// see strings.NewReader, bytes.NewREader & bytes.Buffer for in memory reader
	file, err := os.Open("post.go")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	defer file.Close()

	resp, err := http.Post("https://httpbin.org/post", "application/octet-stream", file)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	defer resp.Body.Close()
	// TODO: check status code

	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatalf("error: %s", err)
	}
}
