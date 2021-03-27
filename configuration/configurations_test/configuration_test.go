package configurations_test

import (
	"github.com/franela/goblin"
	"meshed/configuration/configurations"
	"meshed/meshnode/testsupport"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testsupport.ReadFlags()
	os.Exit(m.Run())
}

func TestConfigCreation(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Read database config file", func() {
		dbConfig := configurations.ReadDbConfig("mesh.db.properties.ini")
		g.It("Db properties should be recognised", func() {
			g.Assert(dbConfig.Dbname).Equal("meshtestdb")
		})
	})
	g.Describe("Read image database config file", func() {
		dbConfig := configurations.ReadDbConfig("imagestream.db.properties.ini")
		g.It("Db properties should be recognised", func() {
			g.Assert(dbConfig.Dbname).Equal("imagetestdb")
		})
	})
}
