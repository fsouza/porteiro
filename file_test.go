// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package porteiro

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestOpenFiles(t *testing.T) {
	t.Parallel()
	opener := OpenFiles(nil)
	abs, err := filepath.Abs("testdata/somefile.txt")
	if err != nil {
		t.Fatal(err)
	}
	uris := []string{
		"testdata/somefile.txt",
		"file://testdata/somefile.txt",
		abs,
		"file://" + abs,
	}
	for _, uri := range uris {
		uri := uri
		t.Run(uri, func(t *testing.T) {
			t.Parallel()
			rc, err := opener.Open(uri)
			if err != nil {
				t.Fatal(err)
			}
			defer rc.Close()
			data, err := ioutil.ReadAll(rc)
			if err != nil {
				t.Fatal(err)
			}
			expectedData := "hello, it's me\n"
			if string(data) != expectedData {
				t.Errorf("wrong data returned\nwant %q\ngot  %q", expectedData, string(data))
			}
		})
	}
}
