package mongodb

import (
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
			configFileName:     "../config/mesh.db.properties.ini.sample",
			desc:               "Read database config file",
			configPropertyName: "meshtestdb",
			configPropertyURL:  "mongodb://user:password@host:port",
		},
	}
	It("should be read the correct configuration without errors", func() {
		for _, tc := range testCases {
			dbconfig, err := decodeDbConfig(tc.configFileName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(dbconfig.MeshDbName).Should(Equal(tc.configPropertyName))
		}
	})
})
