package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"fmt"
)

var sPoint = image.Point{0, 0}
var currentFilePath = "./sourcedata/"

type Sector struct {
	Min, Max int
}

func setSector(para []string) []Sector {

	var tempSector Sector
	var sectors []Sector

	for _, v := range para {

		tempString := strings.Split(v, "-")

		tempSector.Min, _ = strconv.Atoi(tempString[0])
		tempSector.Max, _ = strconv.Atoi(tempString[1])

		sectors = append(sectors, tempSector)
	}

	return sectors
}

func drawImage(fileName string, fixedSectors []Sector) {
	filePath := currentFilePath + fileName
	for i, v := range fixedSectors {
		f, err := os.Open(filePath)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		// fmt.Println("The image has been successfully opened")
		img, _, err := image.Decode(f)

		if err != nil {
			log.Panic(err)
		}

		opt := jpeg.Options{
			Quality: 100,
		}

		startingPoint := image.Point{img.Bounds().Min.X, v.Min}
		cuttingPoint := image.Point{img.Bounds().Max.X, v.Max}

		imageSize := image.Rectangle{startingPoint, cuttingPoint}

		newImgRGBA := image.NewRGBA(image.Rectangle{startingPoint, cuttingPoint})

		draw.Draw(newImgRGBA, imageSize, img, startingPoint, draw.Src)

		tempfileName := strings.Split(fileName, ".")
		divideImageFileName := tempfileName[0] + "_" + strconv.Itoa(i) + ".jpeg"
		dataPath := "./resultdata/" + divideImageFileName
		newImgFile, err := os.Create(dataPath)

		jpeg.Encode(newImgFile, newImgRGBA, &opt)

	}

}

func main() {

	fmt.Println("Image Cutter Program has begin")

	args := os.Args[:]

	files, err := ioutil.ReadDir("./sourcedata")
	if err != nil {
		log.Panic(err)
	}

	var fileList []string

	for _, f := range files {

		fileCheck := strings.Split(f.Name(), ".")

		if len(fileCheck) == 2 {

			if fileCheck[1] == "jpeg" || fileCheck[1] == "jpg" {
				// fmt.Println(f.Name())
				fileList = append(fileList, f.Name())
			}
		}
	}

	fixedSectors := setSector(args[1:])

	for _, fileName := range fileList {
		// fmt.Print("\033[G\033[K") //restore the cursor position and clear the line
		fmt.Printf("Cutting %s into pieces please wait... \n", fileName)
		drawImage(fileName, fixedSectors)
		// fmt.Print("\033[A") // move the cursor up
	}

}
