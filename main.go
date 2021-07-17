package main

import (
	"fmt"
	"time"

	"github.com/kataras/iris/v12"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"

	_ "github.com/frank-bee/go-web-test/docs"
)

type payload struct {
	Message string `json:"message" xml:"message" msgpack:"message" yaml:"Message" url:"message" form:"message"`
}

type User struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Firstname  string        `json:"firstname"`
	Lastname   string        `json:"lastname"`
	Age        int           `json:"age"`
	Email      string        `json:"email"`
	InsertedAt time.Time     `json:"inserted_at" bson:"inserted_at"`
	LastUpdate time.Time     `json:"last_update" bson:"last_update"`
}

var c *mgo.Collection

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	app := iris.Default()
	app.Use(myMiddleware)

	//DB setup
	session, err := mgo.Dial("127.0.0.1")
	if nil != err {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c = session.DB("usergo").C("profiles")

	// Index
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	//swagger
	config := &swagger.Config{
		URL: "http://localhost:8080/swagger/doc.json", //The url pointing to API definition
	}
	// use swagger middleware to
	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))

	// My Endpoints
	app.Handle("GET", "/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "pong"})
	})

	app.Handle("GET", "/users", getUsers)
	app.Handle("GET", "/users/{email}", getUser)
	app.Handle("POST", "/users", createUser)
	app.Handle("PATCH", "/users/{email}", updateUser)
	app.Handle("DELETE", "/users/{email}", deleteUser)

	// Listens and serves incoming http requests
	// on http://localhost:8080.
	app.Listen(":8080")
}

func myMiddleware(ctx iris.Context) {
	ctx.Application().Logger().Infof("Runs before %s", ctx.Path())
	ctx.Next()
}

// @Summary Get all users
// @Description Get all users
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /users [get]
func getUsers(ctx iris.Context) {
	results := []User{}

	err := c.Find(nil).All(&results)
	if err != nil {
		// TODO: Do something about the error
	} else {
		fmt.Println("Results All: ", results)
	}
	ctx.JSON(iris.Map{"response": results})
}

// @Summary Get user
// @Description this to get a user by email
// @Accept  json
// @Produce json
// @Param   email     path    string     true        "Some ID"
// @Success 200 {string} string	"ok"
// @Router /users/{email} [get]
func getUser(ctx iris.Context) {
	email := ctx.Params().Get("email")
	fmt.Println(email)
	if email == "" {
		ctx.JSON(iris.Map{"response": "please pass a valid email"})
	}
	result := User{}
	err := c.Find(bson.M{"email": email}).One(&result)
	if err != nil {
		ctx.JSON(iris.Map{"response": err.Error()})
	}
	ctx.JSON(iris.Map{"response": result})
}

// @Summary Create user
// @Description this to create a user
// @Accept  json
// @Produce  json
// @Param   persondata     body    string     true        "Person Data"
// @Success 200 {string} string	"ok"
// @Router /users [post]
func createUser(ctx iris.Context) {
	params := &User{}

	err := ctx.ReadJSON(params)
	if err != nil {
		ctx.JSON(iris.Map{"response": err.Error()})
	} else {
		params.LastUpdate = time.Now()
		err := c.Insert(params)
		if err != nil {
			ctx.JSON(iris.Map{"response": err.Error()})
		} else {
			fmt.Println("Successfully inserted into database")
			result := User{}
			err = c.Find(bson.M{"email": params.Email}).One(&result)
			if err != nil {
				ctx.JSON(iris.Map{"response": err.Error()})
			}
			ctx.JSON(iris.Map{"response": "User succesfully created", "message": result})
		}
	}

}

// @Summary Update user
// @Description this to update a user by email
// @Accept  json
// @Produce json
// @Param   email     path    string     true        "Some ID"
// @Success 200 {string} string	"ok"
// @Router /users/{email} [patch]
func updateUser(ctx iris.Context) {
	email := ctx.Params().Get("email")
	fmt.Println(email)
	if email == "" {
		ctx.JSON(iris.Map{"response": "please pass a valid email"})
	}
	params := &User{}
	err := ctx.ReadJSON(params)
	if err != nil {
		ctx.JSON(iris.Map{"response": err.Error()})
	} else {
		params.InsertedAt = time.Now()
		query := bson.M{"email": email}
		err = c.Update(query, params)
		if err != nil {
			ctx.JSON(iris.Map{"response": err.Error()})
		} else {
			result := User{}
			err = c.Find(bson.M{"email": params.Email}).One(&result)
			if err != nil {
				ctx.JSON(iris.Map{"response": err.Error()})
			}
			ctx.JSON(iris.Map{"response": "user record successfully updated", "data": result})
		}
	}

}

// @Summary Delete user
// @Description this to delete a user by email
// @Accept  json
// @Produce  json
// @Param   email     path    string     true        "Some ID"
// @Success 200 {string} string	"ok"
// @Router /users/{email} [delete]
func deleteUser(ctx iris.Context) {
	email := ctx.Params().Get("email")
	fmt.Println(email)
	if email == "" {
		ctx.JSON(iris.Map{"response": "please pass a valid email"})
	}
	params := &User{}
	err := ctx.ReadJSON(params)
	if err != nil {
		ctx.JSON(iris.Map{"response": err.Error()})
	} else {
		params.InsertedAt = time.Now()
		query := bson.M{"email": email}
		err = c.Remove(query)
		if err != nil {
			ctx.JSON(iris.Map{"response": err.Error()})
		} else {
			ctx.JSON(iris.Map{"response": "user record successfully deleted"})
		}
	}
}
