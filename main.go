package main

import (
	"encoding/xml"
	"io/ioutil"
	"math"
	"os"

	"github.com/hauke96/kingpin"
	"github.com/hauke96/sigolo"
)

const VERSION string = "v0.0.1"

var (
	app       = kingpin.New("GPX smoother", "A simple application to smooth GPX-tracks written in go")
	appFile   = app.Arg("file", "The .gpx file to smooth").Required().String()
	appOutput = app.Flag("output", "The output file").Short('o').Default("out.gpx").String()
	appWeight = app.Flag("weight", "Specifies how strong the smoothing should happen. Larger numbers result in a more precise track, lower numbers in a smoother one. Default: 3.0").Short('w').Default("3.0").Float64()
	appSize   = app.Flag("size", "Specifies how much surrounding point of each GPX-point should be considered. Larger numbers result in a more precise track, lower numbers in a smoother one. Default: 6").Short('s').Default("6").Int()
	appDebug  = app.Flag("debug", "Verbose mode, showing additional debug information").Short('d').Bool()
)

func configureCliArgs() {
	app.Author("Hauke Stieler")
	app.Version(VERSION)

	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')
}

func configureLogging() {
	if *appDebug {
		sigolo.LogLevel = sigolo.LOG_DEBUG
	} else {
		sigolo.LogLevel = sigolo.LOG_INFO
	}
}

func main() {
	configureCliArgs()

	app.Parse(os.Args[1:])

	configureLogging()

	sigolo.Info("Reading GPX file '%s'", *appFile)

	gpxObj := readGpx(*appFile)

	stepsBackwards := *appSize
	weightFactor := *appWeight

	newTrkpt := make([]Trkpt, len(gpxObj.Trk.Trkseg.Trkpt)-stepsBackwards)

	points := gpxObj.Trk.Trkseg.Trkpt

	sigolo.Info("Start smoothing")

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

		// math.Pow(stepsBackwards+1,2) : The amount of weighted points we need to divide by to get the average
		newTrkpt[i-stepsBackwards/2] = Trkpt{
			Lat:  latSum / weightSum,
			Lon:  lonSum / weightSum,
			Ele:  points[i].Ele,
			Time: points[i].Time,
		}
	}

	sigolo.Info("Finished smoothing")

	newGpx := gpx{}
	newGpx.Xmlns = gpxObj.Xmlns
	newGpx.Version = gpxObj.Version
	newGpx.Trk.Name = gpxObj.Trk.Name
	newGpx.Trk.Trkseg.Trkpt = newTrkpt

	writeGpx(*appOutput, &newGpx)
}

func readGpx(fileName string) *gpx {
	// Open xmlFile
	xmlFile, err := os.Open(fileName)
	sigolo.FatalCheck(err)

	sigolo.Debug("Successfully opened '%s'", fileName)
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	sigolo.FatalCheck(err)

	gpxObj := gpx{}

	err = xml.Unmarshal(byteValue, &gpxObj)
	sigolo.FatalCheck(err)

	sigolo.Debug("Successfully read outputinput GPX-file '%s'", fileName)

	return &gpxObj
}

func writeGpx(fileName string, gpxObj *gpx) {
	xmlFile, err := os.Create(fileName)
	sigolo.FatalCheck(err)

	sigolo.Debug("Successfully opened '%s'", fileName)
	defer xmlFile.Close()

	byteValue, err := xml.MarshalIndent(gpxObj, "", "\t")
	sigolo.FatalCheck(err)

	ioutil.WriteFile(fileName, byteValue, 0644)
	sigolo.FatalCheck(err)

	sigolo.Info("Successfully written output GPX-file '%s'", fileName)
}
