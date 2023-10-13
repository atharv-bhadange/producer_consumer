package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/atharv-bhadange/producer_consumer/database"
	"github.com/atharv-bhadange/producer_consumer/models"
	"github.com/joho/godotenv"
	"github.com/wagslane/go-rabbitmq"
)

func main() {

	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = database.InitDatabase()

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	db := database.DB.Db

	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost:5672/",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Connected to RabbitMQ")

	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))

			id, err := strconv.Atoi(string(d.Body))

			if err != nil {
				log.Println(err)
				return rabbitmq.NackDiscard
			}

			var product models.Product

			if err := db.Where("product_id = ?", id).First(&product).Error; err != nil {
				log.Println(err)
				return rabbitmq.NackDiscard
			}

			path, err := downloadZip(product.ProductImages, product.ProductName)
			pathList := []string{path}

			if err := db.Save(&models.Product{
				ProductID:               id,
				CompressedProductImages: pathList,
				ProductName:             product.ProductName,
				ProductDescription:      product.ProductDescription,
				ProductPrice:            product.ProductPrice,
				ProductImages:           product.ProductImages,
			}).Error; err != nil {
				log.Println(err)
				return rabbitmq.NackDiscard
			}

			if err != nil {
				log.Println(err)
				return rabbitmq.NackDiscard
			}

			return rabbitmq.Ack
		},
		"my_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("secret_string"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close()

	close := make(chan bool)

	log.Println("Waiting for messages...")

	<-close

}

func downloadImage(url string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func downloadZip(imageURLs []string, productName string) (string, error) {

	zipFile, err := os.Create(productName+".zip")
	if err != nil {
		fmt.Println("Error creating ZIP file:", err)
		return "", err
	}
	defer zipFile.Close()

	// Create a new ZIP archive writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Download each image and add it to the ZIP archive
	for i, imageURL := range imageURLs {
		filename := fmt.Sprintf("image%d.jpg", i+1) // Change the filename as needed
		if err := downloadImage(imageURL, filename); err != nil {
			fmt.Printf("Error downloading image %s: %v\n", imageURL, err)
			continue
		}

		fileToZip, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Error opening image file %s: %v\n", filename, err)
			continue
		}
		defer fileToZip.Close()

		// Create a new file header for the ZIP entry
		header := &zip.FileHeader{
			Name:   filename,
			Method: zip.Deflate,
		}

		// Add the image file to the ZIP archive
		fileInZip, err := zipWriter.CreateHeader(header)
		if err != nil {
			fmt.Printf("Error creating ZIP entry for %s: %v\n", filename, err)
			continue
		}

		_, err = io.Copy(fileInZip, fileToZip)
		if err != nil {
			fmt.Printf("Error adding %s to ZIP archive: %v\n", filename, err)
			continue
		}

		if err := os.Remove(filename); err != nil {
			fmt.Printf("Error removing image file %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Added %s to ZIP archive.\n", filename)
	}

	// Return the name of the ZIP file with full path

	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return wd + "/" + productName +".zip", nil
}
