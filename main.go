package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Scarf struct {
	Id           int     `json:"id"`
	Material     string  `json:"material"`
	Price        float64 `json:"price"`
	Manufacturer string  `json:"manufacturer"`
	Color        string  `json:"color"`
	Width        int     `json:"width"`
	Length       int     `json:"length"`
	Size         float64 `json:"size"`
}

type Data struct {
	Price float64
	C1    float64
	C2    float64
	C3    float64
	// for calculating distance
	D1 float64
	D2 float64
	D3 float64
	// updated centers
	M1 float64
	M2 float64
	M3 float64
}

type Centroid struct {
	X float64
	Y float64
}

type Cluster struct {
	X        float64
	Y        float64
	Distance float64
}

func main() {
	scarvesMap := make(map[int]*Data)
	rand.Seed(time.Now().UnixNano())
	for idx, s := range Scarves {
		var c1, c2, c3 float64
		c1 = float64((rand.Intn(5) + 1))/ float64(10)
		c2 = float64((rand.Intn(5) + 1))/ float64(10)

		if (c1 + c2) >= 1 {
			c1 -= 0.1
			c2 -= 0.1
		} else {
			c3 = math.Round((1 - (c1 + c2)) * 10) / float64(10)
		}

		scarvesMap[idx + 1] = &Data{
			Price: s.Price,
			C1:    c1,
			C2:    c2,
			C3:    c3,
		}
	}


	centroid1, centroid2, centroid3 := calculateCentroids(scarvesMap)

	findDistance(scarvesMap, centroid1, centroid2, centroid3)

	updateMembership(scarvesMap)

	for idx, data := range scarvesMap {
		fmt.Println(idx, " : ", data)
	}

}

func calculateCentroids(scarvesMap map[int]*Data) (c1, c2, c3 Centroid) {
	var (
		numeratorC1_X, numeratorC1_Y,
		numeratorC2_X, numeratorC2_Y,
		numeratorC3_X, numeratorC3_Y float64
	)
	var denominatorC1, denominatorC2, denominatorC3 float64

	for idx, data := range scarvesMap {
		c1Pow := math.Pow(data.C1, 2)
		numeratorC1_X += float64(idx) * c1Pow
		numeratorC1_Y += data.Price * c1Pow
		denominatorC1 += c1Pow

		c2Pow := math.Pow(data.C2, 2)
		numeratorC2_X += float64(idx) * c2Pow
		numeratorC2_Y += data.Price * c2Pow
		denominatorC2 += c2Pow

		c3Pow := math.Pow(data.C3, 2)
		numeratorC3_X += float64(idx) * c3Pow
		numeratorC3_Y += data.Price * c3Pow
		denominatorC3 += c3Pow
	}

	c1 = Centroid{
		X: math.Round(numeratorC1_X/denominatorC1*100) / 100,
		Y: math.Round(numeratorC1_Y/denominatorC1*100) / 100,
	}

	c2 = Centroid{
		X: math.Round(numeratorC2_X/denominatorC2*100) / 100,
		Y: math.Round(numeratorC2_Y/denominatorC2*100) / 100,
	}

	c3 = Centroid{
		X: math.Round(numeratorC3_X/denominatorC3*100) / 100,
		Y: math.Round(numeratorC3_Y/denominatorC3*100) / 100,
	}

	return
}

func findDistance(scarvesMap map[int]*Data, centroid1, centroid2, centroid3 Centroid) {
	for idx, data := range scarvesMap {
		scarvesMap[idx].D1 = math.Round(math.Sqrt(math.Pow(centroid1.X-float64(idx), 2)+math.Pow(centroid1.Y-data.Price, 2))*100) / 100
		scarvesMap[idx].D2 = math.Round(math.Sqrt(math.Pow(centroid2.X-float64(idx), 2)+math.Pow(centroid2.Y-data.Price, 2))*100) / 100
		scarvesMap[idx].D3 = math.Round(math.Sqrt(math.Pow(centroid3.X-float64(idx), 2)+math.Pow(centroid3.Y-data.Price, 2))*100) / 100
	}
}

func updateMembership(scarvesMap map[int]*Data) {
	for idx, data := range scarvesMap {
		scarvesMap[idx].M1 = math.Round(data.D1 / (data.D1 + data.D2 + data.D3)* 10) / 10
		scarvesMap[idx].M2 = math.Round(data.D2 / (data.D1 + data.D2 + data.D3)* 10) / 10
		scarvesMap[idx].M3 = math.Round(data.D3 / (data.D1 + data.D2 + data.D3)* 10) / 10
	}
}
