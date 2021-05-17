package main

import (
	"fmt"
	"github.com/lhhong/go-fcm/fcm"
	"math"
	"sort"
)

// Your current data structure
type Point struct {
	X float64
	Y float64
}

// To implement Interface containing Multiply, Add and Norm
// Custom operators can be defined for different data types
type FcmPoint Point

// Multiplying a data point with a scalar weight
func (p FcmPoint) Multiply(weight float64) fcm.Interface {
	return FcmPoint{
		X: weight * p.X,
		Y: weight * p.Y,
	}
}

// Adding 2 data points together
func (p FcmPoint) Add(p2I fcm.Interface) fcm.Interface {
	p2 := p2I.(FcmPoint)
	return FcmPoint{
		X: p2.X + p.X,
		Y: p2.Y + p.Y,
	}
}

// Evaluating distance measure between 2 data points
func (p FcmPoint) Norm(p2I fcm.Interface) float64 {
	p2 := p2I.(FcmPoint)
	xDiff := p.X - p2.X
	yDiff := p.Y - p2.Y
	return math.Sqrt(math.Pow(xDiff, 2.0) + math.Pow(yDiff, 2.0))
}

func FindFCM(fcmPoints []fcm.Interface, value string) map[int]float64 {

	clusteredData := make(map[int]float64, len(fcmPoints))
	// Retrieve centroids and weights of each (centroid, data point) pair
	centroids, weights := fcm.Cluster(fcmPoints, 2.0, 0.00001, 3)

	var centers []float64
	var min, avg, max float64

	for idx := range centroids {
		centers = append(centers, centroids[idx].(FcmPoint).Y)
	}

	sort.Float64s(centers)
	if len(centers) == 3 {
		min, avg, max = centers[0], centers[1], centers[2]
	} else {
		// return empty slice
		return clusteredData
	}

	for i, r := range weights {
		for j, w := range r {
			fmt.Printf("Centroid (%f, %f), Element (%f, %f), weight %f\n",
				centroids[i].(FcmPoint).X, centroids[i].(FcmPoint).Y,
				fcmPoints[j].(FcmPoint).X, fcmPoints[j].(FcmPoint).Y, w)

			if value == "min" && min == centroids[i].(FcmPoint).Y {
				clusteredData[int(fcmPoints[j].(FcmPoint).X)] = w
			} else if value == "avg" && avg  == centroids[i].(FcmPoint).Y {
				clusteredData[int(fcmPoints[j].(FcmPoint).X)] = w
			} else if value == "max" && max == centroids[i].(FcmPoint).Y {
				clusteredData[int(fcmPoints[j].(FcmPoint).X)] = w
			}
		}

		if len(clusteredData) == len(fcmPoints) {
			return clusteredData
		}
	}

	return clusteredData
}
