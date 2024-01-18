package img

import (
	"fmt"
	"image"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

func ConvertImage(c *fiber.Ctx) error {

	// Get image from request
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("failed to get image from request: %v", err)})
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to open file: %v", err)})
	}

	mimeType := fileHeader.Header.Get("content-type")

	switch {
	case strings.Contains(mimeType, "image/heic"), strings.Contains(mimeType, "image/avif"), strings.Contains(mimeType, "image/heif"):
		img, _, err := image.Decode(file)
		if err != nil {
			// The file couldn't be decoded
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to convert file: %v", err)})
		}

		// Resize the image
		src := imaging.Resize(img, Width, 0, imaging.Lanczos)

		// Set content type
		c.Set(fiber.HeaderContentType, "image/png")

		// Send image back to browser
		err = imaging.Encode(c, src, imaging.PNG)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to encode compressed image: %v", err)})
		}
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to convert file: %v", err)})
	}

	return nil
}
