// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type fakeS3 struct {
	s3.S3
	err      error
	objData  []byte
	openObjs []string
}

func (s *fakeS3) GetObjectRequest(input *s3.GetObjectInput) s3.GetObjectRequest {
	return s3.GetObjectRequest{
		Input:   input,
		Request: s.getAWSRequest(input),
	}
}

func (s *fakeS3) getAWSRequest(input *s3.GetObjectInput) *aws.Request {
	var handlers aws.HandlerList
	handlers.PushBack(func(r *aws.Request) {
		obj := "s3://" + aws.StringValue(input.Bucket) + "/" + aws.StringValue(input.Key)
		s.openObjs = append(s.openObjs, obj)
		reader := bytes.NewReader(s.objData)
		r.Data = &s3.GetObjectOutput{
			Body: ioutil.NopCloser(reader),
		}
		r.Error = s.err
	})
	return &aws.Request{
		Handlers: aws.Handlers{
			Send: handlers,
		},
	}
}
