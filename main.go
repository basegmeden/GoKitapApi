package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/basegmeden/goegitim/models"
	"github.com/basegmeden/goegitim/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type Book struct {
	Yazar   string `json:"yazar"`
	Adi     string `json:"adi"`
	Yayinci string `json:"yayinci"`
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "istek başarisiz"})
		return err
	}

	// Book modelini models.Books'a dönüştür
	bookModel := models.Books{
		Yazar:   &book.Yazar,
		Adi:     &book.Adi,
		Yayinci: &book.Yayinci,
	}

	err = r.DB.Create(&bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Kitap olusturulamadi"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "kitap olusturuldu",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id param bos olamaz",
		})
		return nil
	}

	err := r.DB.Delete(&bookModel, id).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "kayit silinemedi",
		})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Kayit Silindi",
	})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "kitaplar gelmedi"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "kitaplar basariyla getirildi",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id bos",
		})
		return nil
	}

	fmt.Println("id", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "kitap bulunamadi"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "id eslesti",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Database yüklenemedi.")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("db migrate edilemedi")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")

	log.Println("Server 8080 portunda başladi...")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
