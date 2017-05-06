// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package porteiro_test

import (
	"fmt"
	"io/ioutil"

	"github.com/fsouza/porteiro"
)

func ExampleOpenFiles_scheme() {
	opener := porteiro.OpenFiles(nil)
	file, err := opener.Open("file://testdata/somefile.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", data)
	// Output: hello, it's me
}

func ExampleOpenFiles_noscheme() {
	opener := porteiro.OpenFiles(nil)
	file, err := opener.Open("testdata/somefile.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", data)
	// Output: hello, it's me
}
