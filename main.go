package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

func main() {
	gpxObj := readGpx("test.gpx")

	stepsBackwards := 6
	weightFactor := 3.0

	newTrkpt := make([]Trkpt, len(gpxObj.Trk.Trkseg.Trkpt)-stepsBackwards)

	points := gpxObj.Trk.Trkseg.Trkpt

	for i := stepsBackwards / 2; i < len(points)-stepsBackwards/2; i++ {
		latSum := 0.0
		lonSum := 0.0

		weightSum := 0.0

		for j := i - stepsBackwards/2; j <= i+stepsBackwards/2; j++ {
			weight := stepsBackwards/2 - int(math.Abs(float64(j-i))) + 1
			latSum += points[j].Lat * math.Pow(float64(weight), weightFactor)
			lonSum += points[j].Lon * math.Pow(float64(weight), weightFactor)
			weightSum += math.Pow(float64(weight), weightFactor)
		}

		// math.Pow(stepsBackwards+1,2) : The amount of weighted points we need divide by to get the average
		newTrkpt[i-stepsBackwards/2] = Trkpt{
			Lat: latSum / weightSum,
			Lon: lonSum / weightSum,
		}
	}

	newGpx := gpx{}
	newGpx.Xmlns = gpxObj.Xmlns
	newGpx.Version = gpxObj.Version
	newGpx.Trk.Name = gpxObj.Trk.Name
	newGpx.Trk.Trkseg.Trkpt = newTrkpt

	writeGpx("test.new.gpx", &newGpx)
}

func readGpx(fileName string) *gpx {
	// Open xmlFile
	xmlFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Successfully Opened", fileName)
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gpxObj := gpx{}

	err = xml.Unmarshal(byteValue, &gpxObj)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &gpxObj
}

func writeGpx(fileName string, gpxObj *gpx) {
	xmlFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Successfully Opened", fileName)
	defer xmlFile.Close()

	byteValue, err := xml.MarshalIndent(gpxObj, "", "\t")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ioutil.WriteFile(fileName, byteValue, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
