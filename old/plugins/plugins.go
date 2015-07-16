// Plugins need to be explicitly imported in order to
// get compiled in the final executable
// Each kind of plugin list is imported here
package plugins

import (
	_ "github.com/nathanielc/morgoth/detector/list"
	_ "github.com/nathanielc/morgoth/engine/list"
	_ "github.com/nathanielc/morgoth/fitting/list"
	_ "github.com/nathanielc/morgoth/notifier/list"
)
