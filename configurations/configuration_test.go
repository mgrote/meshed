package configurations_test

import (
	"github.com/mgrote/meshed/configurations"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {
	testCases := []struct {
		configFileName     string
		desc               string
		configPropertyName string
		configPropertyURL  string
	}{
		{
			configFileName:     "configurations_test/mesh.db.properties.ini",
			desc:               "Read database config file",
			configPropertyName: "meshtestdb",
			configPropertyURL:  "mongodb://user:password@host:port",
		},
		{
			configFileName:     "configurations_test/imagestream.db.properties.ini",
			desc:               "Read image database config file",
			configPropertyName: "imagetestdb",
			configPropertyURL:  "mongodb://user:password@host:port",
		},
	}
	It("should be read the correct configuration without errors", func() {
		for _, tc := range testCases {
			dbconfig, err := configurations.ReadDbConfig(tc.configFileName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dbconfig.DbName).Should(Equal(tc.configPropertyName))
		}
	})
})
