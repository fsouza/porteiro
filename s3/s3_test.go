// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3

import (
	"errors"
	"io/ioutil"
	"testing"
)

func TestS3Opener(t *testing.T) {
	t.Parallel()
	client := &fakeS3{}
	o, err := newS3Opener(client)
	if err != nil {
		t.Fatal(err)
	}
	if o.client != client {
		t.Errorf("wrong client stored\nwant %#v\ngot  %#v", client, o.client)
	}
}

func TestS3OpenerNilClient(t *testing.T) {
	t.Parallel()
	o, err := newS3Opener(nil)
	if err != nil {
		t.Fatal(err)
	}
	if o.client == nil {
		t.Error("unexpected <nil> client")
	}
}

func TestOpenS3(t *testing.T) {
	t.Parallel()
	givenData := "hello it's me"
	client := &fakeS3{objData: []byte(givenData)}
	opener, err := Open(client, nil)
	if err != nil {
		t.Fatal(err)
	}
	rc, err := opener.Open("s3://somebucket/someobject.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != givenData {
		t.Errorf("wrong data returned\nwant %q\ngot  %q", givenData, string(data))
	}
}

func TestOpenS3Failure(t *testing.T) {
	t.Parallel()
	prepErr := errors.New("something went wrong")
	client := &fakeS3{err: prepErr}
	opener, err := Open(client, nil)
	if err != nil {
		t.Fatal(err)
	}
	rc, err := opener.Open("s3://somebucket/someobject.txt")
	if err == nil {
		t.Error("unexpected <nil> error when reading the opened resource")
	}
	if rc != nil {
		t.Errorf("unexpected non-nil ReadCloser: %#v", rc)
	}
}
