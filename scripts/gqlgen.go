// AL201204 This script is to work around the incompatibility of
// gqlgen and the Go vendor system.
// It is based on this resolution:
//   https://github.com/99designs/gqlgen/issues/800#issuecomment-514552917

// +build tools

package main

import "github.com/99designs/gqlgen/cmd"

func main() {
	cmd.Execute()
}
