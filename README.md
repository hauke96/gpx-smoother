# gpx-smoother
A simple approach to smooth a GPX track without loosing too much accuracy.

# Result
![The track becomes smoother using a weighted average of the track-point locations](result.png)

# Usage
`usage: GPX smoother [<flags>] <file>`

Soe for example: `gpx-smoother -d -s 10 -w 6 2019-09-19_12-12-16_new.gpx`

## Parameter
Detailed information with flag `-h, --help`.

| parameter | default | description |
|:--:|:--:|:--|
| `-w, --weight` | `3.0` | Specifies how strong the smoothing should happen. Larger numbers result in a more precise track, lower numbers in a smoother one |
| `-s, --size` | `6` | Specifies how much surrounding point of each GPX-point should be considered. Larger numbers result in a more precise track, lower numbers in a smoother one |
| `-d, --debug` | *(not set)* | Verbose mode, showing additional debug information |
| `-h, --help` | *(not set)* | Shows information about the usage |

## What parameter to choose
It depends on the density of the track points.

For tracks with points **every 5 meter** the default parameters (`-w 3 -s 6`) are fine. For dense tracks a larger size-value (e.g. `-s 10`) might be better.

# Algorithm
Currently this takes each point of the track and gets `n` neighbors around it. Then a weighted average of their locations is calculated: Points near by are weighted more than points far away.

```go
 // Basically: How strong is the weight. This is later an exponent.
weightFactor := 3.0

// for each point:

	// for each neighbor of the current point

		// Weight look like: 1 2 3 [4] 3 2 1 (where [4] is the point currently looking at)
		// This is used as basis for "weight^weightFactor"
		weight := stepsBackwards/2 - int(math.Abs(float64(j-i))) + 1

		// sum up the coordinates
		latSum += points[j].Lat * math.Pow(float64(weight), weightFactor)
		lonSum += points[j].Lon * math.Pow(float64(weight), weightFactor)

	// To later get the correct location, we need to know how many locations we used
	weightSum += math.Pow(float64(weight), weightFactor)

// Get the average
newTrkpt[i-stepsBackwards/2] = Trkpt{
	Lat: latSum / weightSum,
	Lon: lonSum / weightSum,
}

```
