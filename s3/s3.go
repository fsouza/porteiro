// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package s3 provides an opener for S3 resources declared in the format
// s3://<bucket-name>/<object-key>.
package s3

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fsouza/porteiro"
)

type s3Opener struct {
	client ClientAPI
}

type ClientAPI interface {
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func newS3Opener(client ClientAPI) (*s3Opener, error) {
	if client == nil {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, err
		}
		client = s3.NewFromConfig(cfg)
	}
	return &s3Opener{client: client}, nil
}

func (o *s3Opener) open(url *url.URL) (io.ReadCloser, error) {
	resp, err := o.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(url.Host),
		Key:    aws.String(strings.TrimLeft(url.Path, "/")),
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// Open returns an opener that is able of loading files from S3 via the "s3"
// scheme.
//
// The s3 client might be set to nil, in which case a new client will be
// created using the default credential provider chain (see
// https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html#credentials-default
// for more details).
func Open(client ClientAPI, o *porteiro.Opener) (*porteiro.Opener, error) {
	opener, err := newS3Opener(client)
	if err != nil {
		return nil, err
	}
	return o.Register("s3", opener.open), nil
}
