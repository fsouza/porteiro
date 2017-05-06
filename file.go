// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package porteiro

import (
	"io"
	"net/url"
	"os"
	"path"
)

// OpenFiles returns an Opener that knows how to open files in the local file
// system.
func OpenFiles(o *Opener) *Opener {
	return o.Register("", openFile).Register("file", openFile)
}

func openFile(resource *url.URL) (io.ReadCloser, error) {
	return os.Open(path.Join(resource.Host, resource.Path))
}
