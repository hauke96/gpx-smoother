package main

type gpx struct {
	Xmlns   string `xml:"xmlns,attr"`
	Version string `xml:"version,attr"`
	Trk     Trk    `xml:"trk"`
}

type Trk struct {
	Name   string `xml:"name"`
	Trkseg Trkseg `xml:"trkseg"`
}

type Trkseg struct {
	Trkpt []Trkpt `xml:"trkpt"`
}

type Trkpt struct {
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
	Ele  float32 `xml:"ele"`
	Time string  `xml:"time"`
}
