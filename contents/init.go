// Package contents maintains the file tree from root directory
package contents

import (
	"fmt"
	"os"
)

var _rootDir string = "."

// var _root

func SetRoot(root string) error {
	_rootDir = root
	if stat, err := os.Stat(_rootDir); err != nil {
		return err
	} else if !stat.IsDir() {
		fmt.Errorf("%s: root path is not a directory", _rootDir)
	}

	// cleanTree(_root)
	return nil
}
