// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-cleanhttp"
)

func TestOpenHTTP(t *testing.T) {
	content := "hello, it's me"
	server, reqs := startServer([]byte(content), false)
	defer server.Close()
	opener := OpenHTTP(http.DefaultClient, nil)
	rc, err := opener.Open(server.URL + "/some-file.txt")
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Errorf("wrong data was read\nwant %q\ngot  %q", content, string(data))
	}
	r := <-reqs
	if r.Method != "GET" {
		t.Errorf("wrong request method sent\nwant %q\ngot  %q", "GET", r.Method)
	}
	if r.URL.Path != "/some-file.txt" {
		t.Errorf("wrong request url\nwant %q\ngot  %q", "/some-file.txt", r.URL.String())
	}
}

func TestOpenHTTPNilClient(t *testing.T) {
	content := "hello, it's me"
	server, reqs := startServer([]byte(content), false)
	defer server.Close()
	opener := OpenHTTP(nil, nil)
	rc, err := opener.Open(server.URL + "/some-file.txt")
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Errorf("wrong data was read\nwant %q\ngot  %q", content, string(data))
	}
	r := <-reqs
	if r.Method != "GET" {
		t.Errorf("wrong request method sent\nwant %q\ngot  %q", "GET", r.Method)
	}
	if r.URL.Path != "/some-file.txt" {
		t.Errorf("wrong request url\nwant %q\ngot  %q", "/some-file.txt", r.URL.String())
	}
}

func TestOpenHTTPSecure(t *testing.T) {
	content := "hello, it's me, and I'm secure!"
	server, reqs := startServer([]byte(content), true)
	defer server.Close()
	transport := cleanhttp.DefaultTransport()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	opener := OpenHTTP(&http.Client{Transport: transport}, nil)
	rc, err := opener.Open(server.URL + "/some-secure-file.txt")
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Errorf("wrong data was read\nwant %q\ngot  %q", content, string(data))
	}
	r := <-reqs
	if r.Method != "GET" {
		t.Errorf("wrong request method sent\nwant %q\ngot  %q", "GET", r.Method)
	}
	if r.URL.Path != "/some-secure-file.txt" {
		t.Errorf("wrong request url\nwant %q\ngot  %q", "/some-file.txt", r.URL.String())
	}
}

func TestOpenHTTPFailure(t *testing.T) {
	opener := OpenHTTP(nil, nil)
	rc, err := opener.Open("http://192.0.2.14:139193/some-file.txt")
	if err == nil {
		t.Error("unexpected <nil> error")
	}
	if rc != nil {
		t.Errorf("unexpected non-nil ReadCloser: %#v", rc)
	}
}

func startServer(content []byte, tls bool) (*httptest.Server, <-chan *http.Request) {
	fn := httptest.NewServer
	if tls {
		fn = httptest.NewTLSServer
	}
	reqs := make(chan *http.Request, 1)
	return fn(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqs <- r
		w.Write(content)
	})), reqs
}
