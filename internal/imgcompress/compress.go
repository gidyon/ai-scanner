package imgcompress

import (
	"fmt"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

func CompressImage(c *fiber.Ctx) error {
	// Get image from request
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("failed to get image from request: %v", err)})
	}

	fmt.Println("Gotten image header")

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to open file: %v", err)})
	}

	fmt.Println("Gotten image")

	// Compress image
	img, err := imaging.Decode(file)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to get image from request: %v", err)})
	}

	fmt.Println("Image decode done")

	// Resize the cropped image to width = 200px preserving the aspect ratio.
	src := imaging.Resize(img, 500, 0, imaging.Lanczos)

	fmt.Println("Compression done")

	c.Set(fiber.HeaderContentType, "image/png")

	// Send image back to browser
	err = imaging.Encode(c, src, imaging.PNG)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to get image from request: %v", err)})
	}

	fmt.Println("sending done")

	// Send back compressed image
	return nil
}
