package main

import (
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Code  string `json:"code" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Price string `json:"price" binding:"required"`
}

var DB *gorm.DB
var products []Product

func FindProducts(c *gin.Context) {
	var products []Product
	DB.Find(&products)

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func FindProduct(c *gin.Context) {
	var product Product
	if err := DB.Where("code = ?", c.Param("code")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func CreateProduct(c *gin.Context) {
	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := Product{Code: input.Code, Name: input.Name, Price: input.Price}
	DB.Create(&product)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func UpdateProduct(c *gin.Context) {
	var product Product
	if err := DB.Where("code = ?", c.Param("code")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Model(&product).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func DeleteProduct(c *gin.Context) {
	var product Product
	if err := DB.Where("code = ?", c.Param("code")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	DB.Delete(&product)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func main() {
	// Dtabase connection
	db, err := gorm.Open("postgres", "user=postgres password=sandy5353 dbname=text port=5432 sslmode=disable")

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Product{})

	DB = db

	r := gin.Default()

	// Routes
	r.GET("/products", FindProducts)
	r.GET("/products/:code", FindProduct)
	r.POST("/products", CreateProduct)
	r.PUT("/products/:code", UpdateProduct)
	r.DELETE("/products/:code", DeleteProduct)

	r.Run()
}
