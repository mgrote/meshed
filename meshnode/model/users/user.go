package users

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"meshed/meshnode"
	"meshed/meshnode/mesh"
	"meshed/meshnode/model"
	"meshed/meshnode/model/categories"
)

const ClassName = "user"

func UserNodeType() mesh.NodeType {
	return meshnode.NewNodeType([]string{"image", categories.ClassName}, ClassName)
}

type User struct {
	Name		string	`json:"name"`
	Forename	string	`json:"forename"`
	Password	string	`json:"password"`
}

// Registers a method to create this node during deserialisation
func init() {
	model.RegisterType("user", func() mesh.MeshNode {
		return meshnode.NewNodeWithContent(UserNodeType(), User{})
	})
}

func NewNode(name string, forename string) mesh.MeshNode {
	user := User{
		Name: name,
		Forename: forename,
	}
	node := meshnode.NewNodeWithContent(UserNodeType(), user)
	node.Save()
	return node
}

func (u *User) SetPassword(pwd string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	u.Password = string(hash)
	log.Println(u.Password)
}

func (u *User) IsPassword(pwd string) bool {

	userpw := []byte(u.Password)
	compare := []byte(pwd)
	err := bcrypt.CompareHashAndPassword(userpw, compare)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GetUser(m mesh.MeshNode) User {
	user, ok := m.GetContent().(User)
	if !ok {
		log.Fatal("could not convert content from ", m)
	}
	return user
}




