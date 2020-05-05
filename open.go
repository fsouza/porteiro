// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// Register returns a new opener with the given scheme and function registered.
func Register(scheme string, fn OpenFunc) *Opener {
	var o *Opener
	return o.Register(scheme, fn)
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

// Merge returns a new opener that combines the functions from both Openers.
//
// In case of scheme conflict, functions from other are chosen.
func (o *Opener) Merge(other *Opener) *Opener {
	registry := make(map[string]OpenFunc)
	if o != nil {
		for k, v := range o.registry {
			registry[k] = v
		}
	}
	for k, v := range other.registry {
		registry[k] = v
	}
	return &Opener{registry: registry}
}

// Open opens the given uri using one of the underlying OpenFunc.
func (o *Opener) Open(uri string) (io.ReadCloser, error) {
	resource, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: make sure it's a valid URI. Parsing error: %w", uri, err)
	}
	fn, ok := o.registry[resource.Scheme]
	if !ok {
		return nil, &UnknownSchemeError{Scheme: resource.Scheme, URI: uri}
	}
	return fn(resource)
}

// UnknownSchemeError represents an error in opening the given URI becaue the
// provided scheme is unknown to the opener.
type UnknownSchemeError struct {
	URI    string
	Scheme string
}

func (err *UnknownSchemeError) Error() string {
	return fmt.Sprintf("can't open %q: unknown scheme %q", err.URI, err.Scheme)
}
