// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gcs provides an opener for Google Cloud Storage (GCS) resources
// declared in the format gcs://<bucket-name>/<object-key>.
package gcs

import (
	"io"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/fsouza/porteiro"
	"golang.org/x/net/context"
)

type gcsOpener struct {
	client *storage.Client
}

func newGCSOpener(client *storage.Client) (*gcsOpener, error) {
	if client == nil {
		var err error
		client, err = storage.NewClient(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return &gcsOpener{client: client}, nil
}

func (o *gcsOpener) open(url *url.URL) (io.ReadCloser, error) {
	object := o.client.Bucket(url.Host).Object(strings.TrimLeft(url.Path, "/"))
	return object.NewReader(context.Background())
}

// Open returns an opener that is able of loading files from Google Cloud
// Storage (GCS) via the "gcs" scheme.
//
// The gcs client might be set to nil, in which case a new client will be
// created using Google's application default credentials (see
// https://developers.google.com/identity/protocols/application-default-credentials
// for more details).
func Open(client *storage.Client, o *porteiro.Opener) (*porteiro.Opener, error) {
	opener, err := newGCSOpener(client)
	if err != nil {
		return nil, err
	}
	return o.Register("gcs", opener.open).Register("gs", opener.open), nil
}
