// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package porteiro provides a common interface with the method for opening
// remote and local resources. It associates schemes with opening functions.
//
// This library comes with one opening function that can be registered for
// opening local files.
//
// Subpackages in the repository include implementation for some common
// resources such as Google Cloud Storage (GCS), Amazon's Simple Storage
// Service (S3) and HTTP(s).
package porteiro
