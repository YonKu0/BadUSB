package main

import (
	"fmt"
	"image"
	"os"
	"log"
	"path/filepath"

	"gocv.io/x/gocv"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
)

// loadTemplates function equivalent to Python's load_templates
func loadTemplates(templatePaths []string) ([]gocv.Mat, error) {
	var templates []gocv.Mat
	for _, path := range templatePaths {
		img := gocv.IMRead(path, gocv.IMReadGrayScale)
		if img.Empty() {
			return nil, fmt.Errorf("failed to load image from path: %s", path)
		}
		templates = append(templates, img)
	}
	return templates, nil
}

// findButtonOnAllScreens function, equivalent to Python's find_button_on_screen for all screens
func findButtonOnAllScreens(templates []gocv.Mat, scaleRange [2]float64, scaleStep int, threshold float32) (image.Point, bool, error) {
	// Get the number of available screens
	numScreens := screenshot.NumActiveDisplays()
	fmt.Printf("%d\n", numScreens)

	// Loop through all screens
	for i := 0; i < numScreens; i++ {
		// Capture the screen
		bounds := screenshot.GetDisplayBounds(i)
		screenImg, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return image.Point{}, false, err
		}

		// Calculate the offset for the current screen
		offsetX := bounds.Min.X
		offsetY := bounds.Min.Y

		// Convert captured image to Mat
		img, err := gocv.ImageToMatRGB(screenImg)
		if err != nil {
			return image.Point{}, false, err
		}
		defer img.Close()

		gocv.CvtColor(img, &img, gocv.ColorRGBToGray)

		var bestMatch float32
		var bestButtonCenter image.Point
		found := false

		for _, template := range templates {
			for scale := scaleRange[0]; scale <= scaleRange[1]; scale += (scaleRange[1] - scaleRange[0]) / float64(scaleStep) {
				resized := gocv.NewMat()
				defer resized.Close()

				gocv.Resize(template, &resized, image.Point{}, scale, scale, gocv.InterpolationDefault)

				result := gocv.NewMat()
				defer result.Close()

				gocv.MatchTemplate(img, resized, &result, gocv.TmCcoeffNormed, gocv.NewMat())
				_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

				if maxVal > threshold {
					center := image.Point{
						X: maxLoc.X + resized.Cols()/2 + offsetX, // Adjusted for screen offset
						Y: maxLoc.Y + resized.Rows()/2 + offsetY, // Adjusted for screen offset
					}

					if !found || maxVal > bestMatch {
						bestMatch = maxVal
						bestButtonCenter = center
						found = true
					}
				}
			}
		}

		if found {
			return bestButtonCenter, true, nil
		}
	}

	return image.Point{}, false, nil
}


// resourcePath function, equivalent to Python's resource_path
func resourcePath(relativePath string) string {
	// Assuming the images are in a folder named "resources" in the same directory as the executable
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, "./allow_button_images", relativePath)
}

func main() {
	// Set template paths
	templatePaths := []string{
		resourcePath("allow_button_image.jpg"),
		resourcePath("allow_button_image1.jpg"),
		resourcePath("allow_button_image2.jpg"),
		resourcePath("allow_button_image3.jpg"),
		resourcePath("allow_button_image4.jpg"),
		resourcePath("allow_button_image5.jpg"),
		// Add other paths
	}

	// Load templates
	templates, err := loadTemplates(templatePaths)
	if err != nil {
		fmt.Printf("Error loading templates: %v\n", err)
		os.Exit(1)
	}

	// Declare bestMatch here
	var bestMatch float32

	// Find button position on all screens
	buttonPosition, found, err := findButtonOnAllScreens(templates, [2]float64{0.5, 1.0}, 5, 0.9)
	if err != nil {
		fmt.Printf("Error finding button on screen: %v\n", err)
		os.Exit(1)
	}

	if found {
		// Print the button position and matching score
		fmt.Printf("Button Position: X=%d, Y=%d, Matching Score: %f\n", buttonPosition.X, buttonPosition.Y, bestMatch)

		// Define fixed pixel offsets
		offsetX := 10 // Adjust this value as needed
		offsetY := 10 // Adjust this value as needed

		// Move mouse to the button position with the fixed pixel offset and add a delay
		fmt.Printf("Moving mouse to X=%d, Y=%d\n", buttonPosition.X+offsetX, buttonPosition.Y+offsetY)
		robotgo.MoveMouse(buttonPosition.X+offsetX, buttonPosition.Y+offsetY)

		// Add a delay before clicking
		robotgo.MilliSleep(300) // 1 second delay

		// Click the left mouse button
		fmt.Println("Clicking the left mouse button")
		// robotgo.Click("left", false)
	} else {
		fmt.Println("Button not found")
	}
}
