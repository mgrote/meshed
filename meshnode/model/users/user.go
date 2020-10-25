package users

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"meshed/meshnode"
	"meshed/meshnode/dbclient"
	"meshed/meshnode/mesh"
	"meshed/meshnode/model"
)

const ClassName = "user"

func UserNodeType() mesh.NodeType {
	return meshnode.NewNodeType([]string{"image", "category"}, ClassName)
}

type User struct {
	Name		string	`json:"name"`
	Forename	string	`json:"forename"`
	Email 		string  `json:"email"`
	LoginName	string	`json:"login"`
	Password	string	`json:"password"`
}

// Registers a method to create this node during deserialisation
func init() {
	log.Println("user init called")
	model.RegisterTypeConverter("user",
		func() *mesh.MeshNode {
		node := meshnode.NewNodeWithContent(UserNodeType(), User{})
		return &node
	})
	model.RegisterContentConverter(ClassName, GetFromMap)
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

func NewNodeFromRegistration(login string, email string, password string) mesh.MeshNode {
	user := User{
		LoginName: login,
		Email: email,
	}
	user.SetPassword(password)
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
		log.Fatal("could not convert content user from ", m)
	}
	return user
}

func GetFromMap(docmap map[string]interface{}) interface{} {
	return User {
		Name:     docmap["name"].(string),
		Forename: docmap["forename"].(string),
		Password: docmap["password"].(string),
	}
}

func FindUserByLogin(login string) (mesh.MeshNode, bool) {
	category , err := dbclient.FindOneByProperty(ClassName, "login", login)
	if err != nil && err.Error() == dbclient.ErrorDocumentNotFound {
		return nil, false
	}
	return category, true
}





