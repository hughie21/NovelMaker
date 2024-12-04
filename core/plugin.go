// Description: This package is used to define the interface of the plugin.
// Author: Hughie21
// Date: 2024-11-21
// license that can be found in the LICENSE file.

package core

import (
	"context"
)

type (
	Pluginer interface {
		Run(context.Context, ...interface{}) (interface{}, error)
		Info() PluginInfo
	}
	PluginInfo struct {
		Name     string
		Type     string
		Priority int
		Version  string
		Author   string
		Email    string
		Source   string
	}
)

func (p *PluginInfo) Info() PluginInfo {
	return *p
}
