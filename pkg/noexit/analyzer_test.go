package noexit_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/a-x-a/go-metric/pkg/noexit"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), noexit.Analyzer, "./...")
}
