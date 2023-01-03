package main

import (
	"analyzer/check"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(check.Analyzer)
}
