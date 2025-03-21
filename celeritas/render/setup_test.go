package render

import (
	"github.com/CloudyKit/jet/v6"
	"os"
	"testing"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./testdata/views"),
	jet.InDevelopmentMode(),
)

var testRenderer = Render{
	Renderer: "",
	RootPath: "",
	JetViews: views,
}

// TestMain anytime tests in this directory runs Go will run this function first
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
