package FacialRecognition

import (
	"fmt"
	"github.com/Kagami/go-face"
	"log"
	"os"
	"path/filepath"
)

// Path to directory with models and test images. Here it's assumed it
// points to the <https://github.com/Kagami/go-face-testdata> clone.
const dataDir = "./assets/go-face-testdata-master"

var (
	modelsDir = filepath.Join(dataDir, "models")
	imagesDir = filepath.Join(dataDir, "images")
)

// This example shows the basic usage of the package: create an
// recognizer, recognize faces, classify them using few known ones.
func GetFacialRecognizer() *face.Recognizer {
	// Init the recognizer.
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	//defer rec.Close()

	files, err := os.ReadDir(imagesDir)
	if err != nil {
		log.Fatalf("Cannot find images directory: %v", err)
	}

	faces := []face.Face{}

	for _, file := range files {
		image := filepath.Join(imagesDir, file.Name())
		face, err := rec.RecognizeSingleFile(image)
		if err != nil {
			fmt.Printf("Can't recognize: %v", err)
			continue
		} else if face == nil {
			fmt.Printf("Face is nil somehow\n")
			continue
		}

		faces = append(faces, *face)
	}

	// Fill known samples. In the real world you would use a lot of images
	// for each person to get better classification results but in our
	// example we just get them from one big image.
	var samples []face.Descriptor
	var cats []int32

	for _, f := range faces {
		samples = append(samples, f.Descriptor)
		// Each face is unique on that image so goes to its own category.
		cats = append(cats, int32(0))
	}
	// Name the categories, i.e. people on the image.
	//labels := []string{
	//	"Jack",
	//}
	// Pass samples to the recognizer.
	rec.SetSamples(samples, cats)

	// Now let's try to classify some not yet known image.
	testImageJack := filepath.Join(dataDir, "Jack.JPEG")
	jack, err := rec.RecognizeSingleFile(testImageJack)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	if jack == nil {
		log.Fatalf("Not a single face on the image")
	}
	catID := rec.Classify(jack.Descriptor)
	if catID < 0 {
		log.Fatalf("Can't classify")
	}

	return rec
}
