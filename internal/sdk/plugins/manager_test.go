package plugins_test

import (
	"path"
	"runtime"
	"testing"

	"github.com/qilin/crm-api/internal/sdk/plugins"
)

var (
	dir string
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir = path.Join(path.Dir(filename), "../../../test/testdata")
}

func TestPluginManager_Load(t *testing.T) {
	t.Log(dir)
	pm := plugins.NewPluginManager()
	e := pm.Load(path.Join(dir, "/plugins/so/plugin.so"))
	if e != nil {
		t.Error(e.Error())
	}
}
