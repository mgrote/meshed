package mesh_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfigurations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mesh Test Suite")
}
