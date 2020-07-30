// Copyright (c) 2020 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package router

import (
	"tailscale.com/types/logger"
)

type RoutingFuncs struct {
	Up   func() error
	Down func() error
	Set  func(*Config) error
}

// funcRouter delegates to the RoutingFuncs it embeds.
// It is useful, for example, to pass configs into ipn-go-bridge on macOS/iOS.
type funcRouter struct {
	funcs RoutingFuncs
}

// NewFuncRouter returns a Router which delegates to the supplied RoutingFuncs.
func NewFuncRouter(logf logger.Logf, funcs RoutingFuncs) (Router, error) {
	return funcRouter{funcs: funcs}, nil
}

func (r funcRouter) Up() error {
	// Bringing up the routes is handled externally.
	if r.funcs.Up != nil {
		return r.funcs.Up()
	}
	return nil
}

func (r funcRouter) Set(cfg *Config) error {
	if cfg == nil {
		cfg = &shutdownConfig
	}
	if r.funcs.Set != nil {
		return r.funcs.Set(cfg)
	}
	return nil
}

func (r funcRouter) Close() error {
	if r.funcs.Set != nil {
		if err := r.funcs.Set(&shutdownConfig); err != nil {
			return err
		}
	}
	if r.funcs.Down != nil {
		return r.funcs.Down()
	}
	return nil
}
