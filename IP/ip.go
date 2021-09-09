package main

// client side
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// func retryRequest(req *http.Request, count int) (*http.Request, error) {

// }

func getIP() (string, error) {
	// have to use https to use the secure socket layer
	// Code Before Timeout
	//resp, err := http.Get("https://httpbin.org/ip")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://httpbin.org/ip", nil)
	//if you want to set headers and what not you can
	req.Header.Set("Accept", "application/json")
	// always have to check for errors, here it's checking if there was an error running the http.Get request. if not, it retruns the error
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// check that our status code is okay, if not, print out status code and message
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %d %s", resp.StatusCode, resp.Status)
	}
	// remeber headers hve to be strings
	// checking to make sure the content we get back is applicationjson
	mimeType := strings.ToLower(resp.Header.Get("Content-Type"))
	if !strings.HasPrefix(mimeType, "application/json") {
		return "", fmt.Errorf("bad return type: %s", mimeType)
	}

	// Annonymous struct -- didn't define as top level struct (so it would need to be defined outside of funcz)
	// if we wanted to use this reply struct in multiple areas we would define in top level to use multiple times
	var reply struct {
		// Origin MUST match the key of what comes back from JSON in HTTP protocol
		//Origin string

		// option two to use IP we then have to set the field to json origin if you want to do your own mapping
		IP string `json:"origin"`
	}
	// limit response to 1mb
	dec := json.NewDecoder(io.LimitReader(resp.Body, 1<<20))
	// checking if when we run decode on the pointer to the reply, there is an error
	if err := dec.Decode(&reply); err != nil {
		return "", nil
	}

	// io.Copy is used for debugging in the beginning to make sure our status code and our content type didn't return any errors
	// in production this obviously wouldn't happen but when writing it in development it's useful
	//io.Copy(os.Stdout, resp.Body)
	return reply.IP, nil

}

func main() {
	ip, err := getIP()
	// the error we get back here is any of the errors from above we return (not print out above) and then we print it below
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Println("IP", ip)
}
