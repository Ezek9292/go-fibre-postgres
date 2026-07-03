package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ezek9292/go-fibre-postgres/models"
	"github.com/ezek9292/go-fibre-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := models.Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "request failed",
		})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "could not create book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(fiber.Map{
		"message": "book created successfully",
		"data":    book,
	})
	return nil
}

func (r *Repository) GetAllBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Book{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "could not get books",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Book{}

	if id == "" {
		context.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "id is required",
		})
		return nil
	}

	err := r.DB.First(bookModel, id).Error
	if err != nil {
		context.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "book not found",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(fiber.Map{
		"message": "book fetched successfully",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) UpdateBook(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Book{}

	if id == "" {
		context.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "id is required",
		})
		return nil
	}

	err := r.DB.First(bookModel, id).Error
	if err != nil {
		context.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "book not found",
		})
		return err
	}

	err = context.BodyParser(bookModel)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "request failed",
		})
		return err
	}

	err = r.DB.Save(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "could not update book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(fiber.Map{
		"message": "book updated successfully",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Book{}

	if id == "" {
		context.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "id is required",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id).Error
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not delete book",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(fiber.Map{
		"message": "book deleted successfully",
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_book", r.CreateBook)
	api.Get("/books", r.GetAllBooks)
	api.Get("/books/:id", r.GetBookById)
	api.Put("/books/:id", r.UpdateBook)
	api.Delete("/books/:id", r.DeleteBook)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	db, err := storage.NewPostgresConnection(config)
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}

	err = db.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatal("Error migrating the database", err)
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
