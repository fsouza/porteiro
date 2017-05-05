// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package porteiro provides a common interface with the method for opening
// remote and local resources. It associates schemes with opening functions.
//
// Subpackages include implementation for some common resources.
package porteiro

import (
	"fmt"
	"io"
	"net/url"
)

// OpenFunc is a function capable of opening a given resource.
type OpenFunc func(uri *url.URL) (io.ReadCloser, error)

// Opener is the base type that can learn how to open things.
//
// It associates the resource's URI scheme with an OpenFunc. Openers are
// immutable.
type Opener struct {
	registry map[string]OpenFunc
}

// Register registers the given function for the given scheme and returns a new
// Opener.
//
// It overrides the scheme in case it has already been registered. Using an
// empty scheme is legal (for opening files, for example).
func (o *Opener) Register(scheme string, fn OpenFunc) *Opener {
	registry := make(map[string]OpenFunc)
	if o != nil {
		for k, v := range o.registry {
			registry[k] = v
		}
	}
	registry[scheme] = fn
	return &Opener{registry: registry}
}

// Open opens the given uri using one of the underlying OpenFunc.
func (o *Opener) Open(uri string) (io.ReadCloser, error) {
	resource, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	fn, ok := o.registry[resource.Scheme]
	if !ok {
		return nil, fmt.Errorf("can't open %q: unknown scheme %q", uri, resource.Scheme)
	}
	return fn(resource)
}
