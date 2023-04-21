package users

import (
	"github.com/mgrote/meshed/mesh"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const TypeName = "user"

func UserNodeType() mesh.NodeType {
	return mesh.NewNodeType([]string{"image", "category"}, TypeName)
}

type User struct {
	Name      string `json:"name"`
	Forename  string `json:"forename"`
	Email     string `json:"email"`
	LoginName string `json:"login"`
	Password  string `json:"password"`
}

// Registers a method to create this node during deserialisation
func init() {
	log.Println("user init called")
	mesh.RegisterTypeConverter("user",
		func() *mesh.Node {
			node := mesh.NewNodeWithContent(UserNodeType(), User{})
			return &node
		})
	mesh.RegisterContentConverter(TypeName, GetFromMap)
}

func NewNode(name string, forename string) mesh.Node {
	user := User{
		Name:     name,
		Forename: forename,
	}
	node := mesh.NewNodeWithContent(UserNodeType(), user)
	node.Save()
	return node
}

func NewNodeFromRegistration(login string, email string, password string) mesh.Node {
	user := User{
		LoginName: login,
		Email:     email,
	}
	user.SetPassword(password)
	node := mesh.NewNodeWithContent(UserNodeType(), user)
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

func GetUser(m mesh.Node) User {
	user, ok := m.GetContent().(User)
	if !ok {
		log.Fatal("could not convert content user from ", m)
	}
	return user
}

func GetFromMap(docmap map[string]interface{}) interface{} {
	return User{
		Name:     docmap["name"].(string),
		Forename: docmap["forename"].(string),
		Password: docmap["password"].(string),
	}
}
