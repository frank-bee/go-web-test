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

type User struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Firstname  string        `json:"firstname"`
	Lastname   string        `json:"lastname"`
	Age        int           `json:"age"`
	Msisdn     string        `json:"msisdn"`
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
		Key:        []string{"msisdn"},
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
	app.Handle("GET", "/users/{some_id} [get]", getUser)
	app.Handle("POST", "/users", createUser)
	app.Handle("PATCH", "/users/{some_id} [get]", updateUser)
	app.Handle("DELETE", "/users/{some_id} [get]", deleteUser)

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
// @Description this to get a user by msisdn
// @Accept  json
// @Produce json
// @Param   some_id     path    string     true        "Some ID"
// @Success 200 {string} string	"ok"
// @Router /users/{some_id} [get]
func getUser(ctx iris.Context) {
	msisdn := ctx.Params().Get("msisdn")
	fmt.Println(msisdn)
	if msisdn == "" {
		ctx.JSON(iris.Map{"response": "please pass a valid msisdn"})
	}
	result := User{}
	err := c.Find(bson.M{"msisdn": msisdn}).One(&result)
	if err != nil {
		ctx.JSON(iris.Map{"response": err.Error()})
	}
	ctx.JSON(iris.Map{"response": result})
}

// @Summary Create user
// @Description this to create a user
// @Accept  json
// @Produce  json
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
			err = c.Find(bson.M{"msisdn": params.Msisdn}).One(&result)
			if err != nil {
				ctx.JSON(iris.Map{"response": err.Error()})
			}
			ctx.JSON(iris.Map{"response": "User succesfully created", "message": result})
		}
	}

}

// @Summary Update user
// @Description this to update a user by msisdn
// @Accept  json
// @Produce json
// @Param   some_id     path    string     true        "Some ID"
// @Success 200 {string} string	"ok"
// @Router /users/{some_id} [patch]
func updateUser(ctx iris.Context) {
	msisdn := ctx.Params().Get("msisdn")
	fmt.Println(msisdn)
	if msisdn == "" {
		ctx.JSON(iris.Map{"response": "please pass a valid msisdn"})
	}
	params := &User{}
	err := ctx.ReadJSON(params)
	if err != nil {
		ctx.JSON(iris.Map{"response": err.Error()})
	} else {
		params.InsertedAt = time.Now()
		query := bson.M{"msisdn": msisdn}
		err = c.Update(query, params)
		if err != nil {
			ctx.JSON(iris.Map{"response": err.Error()})
		} else {
			result := User{}
			err = c.Find(bson.M{"msisdn": params.Msisdn}).One(&result)
			if err != nil {
				ctx.JSON(iris.Map{"response": err.Error()})
			}
			ctx.JSON(iris.Map{"response": "user record successfully updated", "data": result})
		}
	}

}

// @Summary Delete user
// @Description this to delete a user by msisdn
// @Accept  json
// @Produce  json
// @Param   some_id     path    string     true        "Some ID"
// @Success 200 {string} string	"ok"
// @Router /users/{some_id} [delete]
func deleteUser(ctx iris.Context) {
	msisdn := ctx.Params().Get("msisdn")
	fmt.Println(msisdn)
	if msisdn == "" {
		ctx.JSON(iris.Map{"response": "please pass a valid msisdn"})
	}
	params := &User{}
	err := ctx.ReadJSON(params)
	if err != nil {
		ctx.JSON(iris.Map{"response": err.Error()})
	} else {
		params.InsertedAt = time.Now()
		query := bson.M{"msisdn": msisdn}
		err = c.Remove(query)
		if err != nil {
			ctx.JSON(iris.Map{"response": err.Error()})
		} else {
			ctx.JSON(iris.Map{"response": "user record successfully deleted"})
		}
	}
}
