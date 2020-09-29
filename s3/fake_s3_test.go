// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type fakeS3 struct {
	err      error
	objData  []byte
	openObjs []string
}

func (s *fakeS3) GetObject(_ context.Context, input *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	obj := "s3://" + *input.Bucket + "/" + *input.Key
	s.openObjs = append(s.openObjs, obj)
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader(s.objData)),
	}, s.err
}
