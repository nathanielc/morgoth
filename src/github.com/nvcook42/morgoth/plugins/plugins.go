// Plugins need to be explicitly imported in order to
// get compiled in the final executable
// Each kind of plugin list is imported here
package plugins

import (
	_ "github.com/nvcook42/morgoth/detector/list"
	_ "github.com/nvcook42/morgoth/engine/list"
	_ "github.com/nvcook42/morgoth/fitting/list"
)
