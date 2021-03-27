package dbclient

const meshDbConfigFile = "mesh.db.properties.ini"
const gridDbConfigFile = "imagestream.db.properties.ini"

func InitDatabase() {
	InitMeshDbClientWithConfig(meshDbConfigFile)
	InitFileStreamDbClientWithConfig(gridDbConfigFile)
}
