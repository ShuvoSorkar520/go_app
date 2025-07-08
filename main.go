package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var DB *gorm.DB

// func ConnectDB() {
// 	dsn := "root:@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"
// 	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Database connection failed: ", err)
// 	}

// 	database.AutoMigrate(&User{})
// 	DB = database
// }

// Database connect function
func ConnectDB() {
	dsn := "root:JsoEoUbmoVdbjVpeWOAyJUQTKAoPymTU@tcp(mainline.proxy.rlwy.net:42773)/railway?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Auto migrate model
	db.AutoMigrate(&User{})

	DB = db
	log.Println("Connected to MySQL successfully!")
}

func main() {
	app := fiber.New()

	ConnectDB()

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("🚀 Go Fiber App with Railway MySQL is running!")
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		var users []User
		DB.Find(&users)
		return c.JSON(users)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		DB.Create(&user)
		return c.JSON(user)
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var user User
		if result := DB.First(&user, id); result.Error != nil {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		DB.Save(&user)
		return c.JSON(user)
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var user User
		if result := DB.First(&user, id); result.Error != nil {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		DB.Delete(&user)
		return c.SendString("Deleted")
	})

	log.Fatal(app.Listen(":8080"))
}
