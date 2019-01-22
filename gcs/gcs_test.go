// Copyright 2019 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gcs

import (
	"io/ioutil"
	"testing"

	"github.com/fsouza/fake-gcs-server/fakestorage"
)

func TestOpen(t *testing.T) {
	const content = "some content"
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			BucketName: "some-bucket",
			Name:       "files/file.txt",
			Content:    []byte(content),
		},
	})
	defer server.Stop()
	opener, err := Open(server.Client(), nil)
	if err != nil {
		t.Fatal(err)
	}

	gsObj, err := opener.Open("gs://some-bucket/files/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer gsObj.Close()
	data, err := ioutil.ReadAll(gsObj)
	if err != nil {
		t.Fatal(err)
	}
	if strData := string(data); strData != content {
		t.Errorf("wrong content read\nwant %q\ngot  %q", content, strData)
	}

	gcsObj, err := opener.Open("gcs://some-bucket/files/file.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer gcsObj.Close()
	data, err = ioutil.ReadAll(gcsObj)
	if err != nil {
		t.Fatal(err)
	}
	if strData := string(data); strData != content {
		t.Errorf("wrong content read\nwant %q\ngot  %q", content, strData)
	}
}
