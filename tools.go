// +build tools

package main

import (
	_ "github.com/Songmu/ghch/cmd/ghch"
	_ "github.com/Songmu/goxz/cmd/goxz"
	_ "github.com/mattn/goveralls"
	_ "github.com/tcnksm/ghr"
	_ "github.com/x-motemen/gobump/cmd/gobump"
	_ "golang.org/x/lint/golint"
)
