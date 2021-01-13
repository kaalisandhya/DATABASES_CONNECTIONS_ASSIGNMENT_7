package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Id   int    `gorm:" AUTO_INCREMENT "` // increment
	Name string `gorm:" size: 255 "`      // string default length is 255, the use of this tag Reset
	Age  int
}

var (
	db  *gorm.DB
	err error
)

func main() {
	// link mysql
	db, err = gorm.Open("mysql", "root:8603mysql@123@tcp(127.0.0.1:3307)/sandhya?parseTime=true")
	if err != nil {
		panic(err)
	} else {

		db.SingularTable(true) // If set to true, `User` default table named` user`, use `TableName` set the table name will not be affected

		// generally do not directly create a table with CreateTable
		// Check the model `User` table exists, otherwise User` create a table for the model`
		if !db.HasTable(&User{}) {
			if err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}).Error; err != nil {
				panic(err)
			}
		}
	}

	// import routes
	Router()
}

func Router() {
	router := gin.Default()
	// path mapping
	router.GET("/user", InitPage)
	router.POST("/user/create", CreateUser)
	router.GET("/user/list", ListUser)
	router.PUT("/user/update/:id", UpdateUser)
	router.GET("/user/find/:id", GetUser)
	router.DELETE("/user/:id", DeleteUser)

	router.Run(":8080")
}

// Each routing function corresponds to a specific operation, enabling user to add, delete, change, operation
func InitPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK!",
	})
}

// Create user
// curl -i -X POST -H "Content-Type: application/json" -d "{ \"name\": \"Vic\", \"age\": 20}" http://localhost:8080/user/create
func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user) // padding data using bindJson

	// db.Create (& user) // Create Object
	// c.JSON (http.StatusOK, & user) // returns the page
	if user.Name != "" && user.Age > 0 {
		db.Create(&user)
		c.JSON(http.StatusOK, gin.H{"success": &user})
	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}
}

// update user
//  http://localhost:8080/user/update/9
func UpdateUser(c *gin.Context) {
	var user User
	id := c.Params.ByName("id")
	err := db.First(&user, id).Error
	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err.Error())
	} else {
		c.BindJSON(&user)
		db.Save(&user)               // commit the changes
		c.JSON(http.StatusOK, &user) // returns the page
	}
}

// list all users
// http://127.0.0.1:8080/user/list
// curl -i http://localhost:8080/user/list
func ListUser(c *gin.Context) {
	var user []User
	db.Find(&user)
	c.JSON(http.StatusOK, &user) // find the limit line before the line
}

// list the single user
// curl -i http://localhost:8080/user/find/18
func GetUser(c *gin.Context) {
	var user User
	id := c.Params.ByName("id")
	err := db.First(&user, id).Error
	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err.Error())
	} else {
		c.JSON(http.StatusOK, &user)
	}
}

// delete users
// curl -i -X DELETE http://localhost:8080/user/1
func DeleteUser(c *gin.Context) {
	var user User
	id := c.Params.ByName("id")
	db.First(&user, id)
	if user.Id != 0 {
		db.Delete(&user)
		c.JSON(http.StatusOK, gin.H{
			"success": "User# " + id + " deleted!",
		})
	} else {
		c.JSON(404, gin.H{
			"error": "User not found",
		})
	}
}
