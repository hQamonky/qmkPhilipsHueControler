package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	// Create a new router
	r := newRouter()
	// Create a mock server
	mockServer := httptest.NewServer(r)

	// Make a GET request to the "hello" route
	resp, err := http.Get(mockServer.URL + "/connect/")
	if err != nil {
		t.Fatal(err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", resp.StatusCode)
	}

	// It isn't wise to test the entire content of the HTML file.
	// Instead, we test that the content-type header is "text/html; charset=utf-8"
	// so that we know that an html file has been served
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}

}

func TestRouterForNonExistentRoute(t *testing.T) {
	// Create a new router
	r := newRouter()
	// Create a mock server
	mockServer := httptest.NewServer(r)

	// Make POST request that should not work
	resp, err := http.Post(mockServer.URL+"/connect/", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check status code, (should be 405 (method not allowed))
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status should be 405, got %d", resp.StatusCode)
	}

	// Check response body (should be empty)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	respString := string(b)
	expected := ""
	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}

}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	// We want to hit the `GET /connect/` route to get the index.html file response
	resp, err := http.Get(mockServer.URL + "/connect/")
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", resp.StatusCode)
	}

	// It isn't wise to test the entire content of the HTML file.
	// Instead, we test that the content-type header is "text/html; charset=utf-8"
	// so that we know that an html file has been served
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}
}
