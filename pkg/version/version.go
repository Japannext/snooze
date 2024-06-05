package version

import (
	"fmt"
)

var (
	Version string
	Commit  string
)

func Print() {
	fmt.Printf("%s-%s", Version, Commit)
}
