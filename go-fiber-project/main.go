package main

import (
	"gofiberproject/models"
	"gofiberproject/storage"
	"log"
	"net/http"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author	string	`json:"author"` 
	Title  string	 `json:"title"`
	Publisher string	`json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error{
	book := Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message":"request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"couldn't create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"book has been added"})
	return nil

}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"messsage":"id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(bookModel,id)
	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"couldn't delete book"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message":"book has been deleted"})
	return nil
	
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	
	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"Couldn't get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "books fetched succesfully",
			"data" : bookModels,
	})	

	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {
	bookModel := &models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"messsage":"id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message":"Couldn't get book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book fetched succesfully",
			"data" : bookModel,
	})	

	return nil
	
}




func (r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/createBooks", r.CreateBook)
	api.Delete("/deleteBooks/:id", r.DeleteBook)
	api.Get("/getBooks/:id", r.GetBookById)
	api.Get("/Books", r.GetBooks)


}


func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("couldn't load the database")
	}
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("couldn't migrate database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")

}  	