package domain

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"meshnode/mesh"
	"meshnode/meshnode"
)

type userNodeType struct {
	acceptTypes	[]string
	class		string
}

func getUserNodeType() mesh.NodeType {
	return userNodeType{
		[]string{"image", "user"},
		"user"}
}

func (t userNodeType) GetClass() string {
	return t.class
}

func (t userNodeType) AcceptTypes() []string {
	return t.acceptTypes
}

func (t userNodeType) IsAccepted(className string) bool {
	for _, slid := range t.acceptTypes {
		if slid == className {
			return true
		}
	}
	return false
}

type User struct {
	Name		string	`json:"name"`
	Forename	string	`json:"forename"`
	Password	string	`json:"password"`
}

func NewUserNode(name string, forename string) mesh.MeshNode {
	user := User{
		Name: name,
		Forename: forename,
	}
	node := meshnode.NewNodeWithContent(getUserNodeType(), user)
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
		log.Fatal("could not convert user from ", m)
	}
	return user
}





