package img

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

func CompressImage(c *fiber.Ctx) error {
	// Get image from request
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("failed to get image from request: %v", err)})
	}

	width := c.FormValue("width")
	height := c.FormValue("height")

	w, err := strconv.Atoi(width)
	if err != nil {
		w = Width
	}

	h, err := strconv.Atoi(height)
	if err != nil {
		h = 0
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to open file: %v", err)})
	}

	// Compress image
	img, err := imaging.Decode(file)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to decode image: %v", err)})
	}

	// Resize the cropped image to  preserving the aspect ratio.
	src := imaging.Resize(img, w, h, imaging.Lanczos)

	c.Set(fiber.HeaderContentType, "image/png")

	// Send image back to browser
	err = imaging.Encode(c, src, imaging.JPEG)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("failed to emcode compressed image: %v", err)})
	}

	// Send back compressed image
	return nil
}
