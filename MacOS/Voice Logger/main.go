package main

import (
    "fmt"
    "image"
    "log"
    "sync"

    "gocv.io/x/gocv"
    "github.com/go-vgo/robotgo"
    "github.com/kbinani/screenshot"
    "github.com/gobuffalo/packr/v2"
)

// Global variable for the packr box
var box *packr.Box

func init() {
    // Initialize the packr box pointing to the subdirectory with your images
    box = packr.New("myBox", "./allow_button_images")
}

func loadTemplates(templateData [][]byte) ([]gocv.Mat, error) {
    var templates []gocv.Mat
    for _, data := range templateData {
        img, err := gocv.IMDecode(data, gocv.IMReadGrayScale)
        if err != nil {
            return nil, fmt.Errorf("failed to decode image data: %v", err)
        }
        templates = append(templates, img)
    }
    return templates, nil
}

// processScreen function processes a single screen in parallel
func processScreen(screenIndex int, templates []gocv.Mat, scaleRange [2]float64, scaleStep int, threshold float32, resultCh chan image.Point) {
    // Capture the screen
    bounds := screenshot.GetDisplayBounds(screenIndex)
    screenImg, err := screenshot.CaptureRect(bounds)
    if err != nil {
        log.Printf("Error capturing screen %d: %v\n", screenIndex, err)
        resultCh <- image.Point{}
        return
    }

    // Calculate the offset for the current screen
    offsetX := bounds.Min.X
    offsetY := bounds.Min.Y

    // Convert captured image to Mat
    img, err := gocv.ImageToMatRGB(screenImg)
    if err != nil {
        log.Printf("Error converting image to Mat for screen %d: %v\n", screenIndex, err)
        resultCh <- image.Point{}
        return
    }
    defer img.Close()

    gocv.CvtColor(img, &img, gocv.ColorRGBToGray)

    var bestButtonCenter image.Point
    found := false
    var maxValGlobal float32

    var wg sync.WaitGroup
    for _, template := range templates {
        wg.Add(1)
        go func(template gocv.Mat) {
            defer wg.Done()
            for scale := scaleRange[0]; scale <= scaleRange[1]; scale += (scaleRange[1] - scaleRange[0]) / float64(scaleStep) {
                resized := gocv.NewMat()
                gocv.Resize(template, &resized, image.Point{}, scale, scale, gocv.InterpolationDefault)

                result := gocv.NewMat()
                gocv.MatchTemplate(img, resized, &result, gocv.TmCcoeffNormed, gocv.NewMat())
                _, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

                if maxVal > threshold && maxVal > maxValGlobal {
                    maxValGlobal = maxVal
                    bestButtonCenter = image.Point{
                        X: maxLoc.X + resized.Cols()/2 + offsetX,
                        Y: maxLoc.Y + resized.Rows()/2 + offsetY,
                    }
                    found = true
                }
            }
        }(template)
    }
    wg.Wait()

    if found {
        resultCh <- bestButtonCenter
    } else {
        resultCh <- image.Point{}
    }
}

// resourcePath function now retrieves the images from the packr box
func resourcePath(relativePath string) ([]byte, error) {
    data, err := box.Find(relativePath)
    if err != nil {
        return nil, fmt.Errorf("failed to find %s in the box: %v", relativePath, err)
    }
    return data, nil
}

func main() {
    // Initialize a slice of byte slices to hold the image data
    var templateData [][]byte

    // Load each image and append its data to the slice
    imgPaths := []string{
        "allow_button_image.jpg",
        "allow_button_image1.jpg",
        "allow_button_image2.jpg",
        "allow_button_image3.jpg",
        "allow_button_image4.jpg",
        "allow_button_image5.jpg",
        "allow_button_image6.jpg",
        "allow_button_image7.jpg",
        "allow_button_image8.jpg",
        "allow_button_image9.jpg",
    }

    for _, imgPath := range imgPaths {
        data, err := resourcePath(imgPath)
        if err != nil {
            log.Fatalf("Error loading image: %v", err)
        }
        templateData = append(templateData, data)
    }

    // Load templates using the image data
    templates, err := loadTemplates(templateData)
    if err != nil {
        log.Fatalf("Error loading templates: %v", err)
    }

    // Define parameters
    scaleRange := [2]float64{0.5, 1.0}
    scaleStep := 5
    threshold := float32(0.9)

    // Get the number of available screens
    numScreens := screenshot.NumActiveDisplays()

    // Channel to receive results
    resultCh := make(chan image.Point, numScreens)

    // Process screens in parallel
    for i := 0; i < numScreens; i++ {
        go processScreen(i, templates, scaleRange, scaleStep, threshold, resultCh)
    }

    // Wait for results from all screens
    found := false
    var bestButtonCenter image.Point

    for i := 0; i < numScreens; i++ {
        result := <-resultCh
        if result != (image.Point{}) {
            bestButtonCenter = result
            found = true
        }
    }

    // Close the result channel
    close(resultCh)

    if found {
        // Define fixed pixel offsets
        offsetX := 5 
        offsetY := 5 

        robotgo.MoveMouse(bestButtonCenter.X+offsetX, bestButtonCenter.Y+offsetY)

        robotgo.MilliSleep(300) 

        robotgo.Click("left", false)
    } else {
        fmt.Println("Button not found")
    }
}
