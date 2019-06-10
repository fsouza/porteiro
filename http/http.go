// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package http provides an opener for HTTP(s) resources declared in the format
// http(s)://<resource>.
package http

import (
	"io"
	"net/http"
	"net/url"

	"github.com/fsouza/porteiro"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

type httpOpener struct {
	client *http.Client
}

func newHTTPOpener(client *http.Client) *httpOpener {
	if client == nil {
		client = cleanhttp.DefaultClient()
	}
	return &httpOpener{client: client}
}

func (o *httpOpener) open(url *url.URL) (io.ReadCloser, error) {
	resp, err := o.client.Get(url.String()) //nolint:bodyclose
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Open returns an opener that is able of loading files from S3 via the "s3"
// scheme.
//
// The http client might be set to nil, in which case a new client will be
// created.
func Open(client *http.Client, o *porteiro.Opener) *porteiro.Opener {
	opener := newHTTPOpener(client)
	return o.Register("http", opener.open).Register("https", opener.open)
}
