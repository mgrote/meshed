package mongodb

const meshDbConfigFile = "mesh.db.properties.ini"
const gridDbConfigFile = "imagestream.db.properties.ini"

func InitDatabase() error {
	err := InitMeshDbClientWithConfig(meshDbConfigFile)
	if err != nil {
		return err
	}
	return InitFileStreamDbClientWithConfig(gridDbConfigFile)
}
