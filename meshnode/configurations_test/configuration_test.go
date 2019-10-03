package configurations_test

import (
	"github.com/franela/goblin"
	"meshed/meshnode/configurations"
	"testing"
)

func TestConfigCreation(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Config creation", func() {
		dbConfig := configurations.ReadConfig("/Users/michaelgrote/etc/gotest/db.properties.ini")
		g.It("Db properties should be recognised", func() {
			g.Assert(dbConfig.Dbname).Equal("meshdb")
		})
	})
}
