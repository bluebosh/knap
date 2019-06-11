package controller

import (
	"github.com/bluebosh/knap/pkg/controller/appengine"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, appengine.Add)
}
