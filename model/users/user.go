package users

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"meshnode/mesh"
	"meshnode/meshnode"
	"meshnode/model"
)

func UserNodeType() mesh.NodeType {
	return meshnode.NewNodeType([]string{"image", "category"}, "user")
}

type User struct {
	Name		string	`json:"name"`
	Forename	string	`json:"forename"`
	Password	string	`json:"password"`
}

func init() {
	model.RegisterType("category", func() mesh.MeshNode {
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
}

func (u *User) IsPassword(pwd string) bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return u.Password == string(hash)
}

func GetUser(m mesh.MeshNode) User {
	user, ok := m.GetContent().(User)
	if !ok {
		log.Fatal("could not convert content from ", m)
	}
	return user
}





