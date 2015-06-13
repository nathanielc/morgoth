// List All detectors that should be compiled into the final
// morgoth executable
package list

import (
	_ "github.com/nathanielc/morgoth/detector/kstest"
	_ "github.com/nathanielc/morgoth/detector/mgof"
	_ "github.com/nathanielc/morgoth/detector/tukey"
)
