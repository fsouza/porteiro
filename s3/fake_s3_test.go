// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type fakeS3 struct {
	s3.S3
	err      error
	objData  []byte
	openObjs []string
}

func (s *fakeS3) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	obj := "s3://" + aws.StringValue(input.Bucket) + "/" + aws.StringValue(input.Key)
	s.openObjs = append(s.openObjs, obj)
	reader := bytes.NewReader(s.objData)
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(reader),
	}, s.err
}
