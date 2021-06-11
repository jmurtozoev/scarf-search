package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lhhong/go-fcm/fcm"
	"log"
	"net/http"
	"sort"
	"strconv"
)

var cache []Scarf

func init()  {
	cache = make([]Scarf, 0)
	if err := ConnectDb(""); err != nil {
		panic(err)
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl")

	router.GET("/", func(c *gin.Context) {
		scarves, err := GetAll()
		if err != nil {
			panic(err)
		}

		results := make([]Result, 0)
		for _, scarf := range scarves {
			results = append(results, Result{Scarf: scarf})
		}

		c.HTML(http.StatusOK, "main", map[string]interface{}{
			"Price": priceOptions,
			"Length": lengthOptions,
			"Width": widthOptions,
			"Rows": results,
		})
	})

	router.GET("/search", filterScarves)
	router.GET("/personalize", personalize)

	log.Fatal(router.Run(":8080"))

}


func filterScarves(c *gin.Context) {
	var scarves []Scarf
	var lengthOption, widthOption, priceOption int
	var filtersCount int
	var priceCluster, widthCluster, lengthCluster map[int]float64
	results := make([]Result, 0)

	scarves, err := GetAll()
	if err != nil {
		log.Println("Error getting scarves list")
	}

	if priceOption, err = strconv.Atoi(c.Query("price")); err == nil && priceOption > 1 {
		value := "min"
		switch priceOption {
		case 3:
			value = "avg"
		case 4:
			value = "max"
		}

		fcmPoints := make([]fcm.Interface, len(scarves))
		for i, s := range scarves {
			fcmPoints[i] = FcmPoint(Point{
				float64(s.Id),
				s.Price,
			})
		}
		priceCluster = FindFCM(fcmPoints, value)
		filtersCount++
	}


	if widthOption, err = strconv.Atoi(c.Query("width")); err == nil && widthOption > 1 {
		value := "min"
		switch widthOption {
		case 3:
			value = "avg"
		case 4:
			value = "max"
		}

		fcmPoints := make([]fcm.Interface, len(scarves))
		for i, s := range scarves {
			fcmPoints[i] = FcmPoint(Point{
				float64(s.Id),
				float64(s.Width),
			})
		}
		widthCluster = FindFCM(fcmPoints, value)
		filtersCount++
	}

	if lengthOption, err = strconv.Atoi(c.Query("length")); err == nil && lengthOption > 1 {
		value := "min"
		switch lengthOption {
		case 3:
			value = "avg"
		case 4:
			value = "max"
		}

		fcmPoints := make([]fcm.Interface, len(scarves))
		for i, s := range scarves {
			fcmPoints[i] = FcmPoint(Point{
				float64(s.Id),
				float64(s.Length),
			})
		}
		lengthCluster = FindFCM(fcmPoints, value)
		filtersCount++
	}

	for _, s := range scarves {
		result := Result{
			Scarf: s,
			TotalWeight: 1,
		}

		if v, ok := priceCluster[s.Id]; ok {
			result.TotalWeight *= v
		}

		if v, ok := lengthCluster[s.Id]; ok {
			result.TotalWeight *= v
		}

		if v, ok := widthCluster[s.Id]; ok {
			result.TotalWeight *= v
		}

		if filtersCount == 1 && result.TotalWeight > 0.5 {
			results = append(results, result)
		} else if filtersCount == 2 && result.TotalWeight > 0.25 {
			results = append(results, result)
		} else if filtersCount == 3 && result.TotalWeight > 0.1 {
			results = append(results, result)
		} else if filtersCount == 0 {
			results = append(results, result)
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i].TotalWeight > results[j].TotalWeight})

	updatedScarves := make([]Scarf, 0)
	for _, r := range results {
		updatedScarves = append(updatedScarves, r.Scarf)
	}

	if len(results) >= 3 {
		c.HTML(http.StatusOK, "personalization", map[string]interface{}{
			"SelectedPrice": priceOption,
			"SelectedLength": lengthOption,
			"SelectedWidth": widthOption,
			"Price": priceOptions,
			"Length": lengthOptions,
			"Width": widthOptions,
			"Rows": results,
		})

		// set sorted scarves into cache
		cache = updatedScarves
		copy(cache, updatedScarves)
		fmt.Println(len(cache), cap(cache), len(updatedScarves))
	} else {
		c.HTML(http.StatusOK, "final", map[string]interface{}{
			"Rows": results,
		})
	}
}

func personalize(c *gin.Context) {
	var scarves []Scarf
	var filtersCount int
	var priceCluster, widthCluster, lengthCluster map[int]float64
	results := make([]Result, 0)

	scarves = cache
	fmt.Println(len(cache))

	if priceOption, err := strconv.Atoi(c.Query("price")); err == nil && priceOption > 1 {
		value := "min"
		switch priceOption {
		case 3:
			value = "avg"
		case 4:
			value = "max"
		}

		fcmPoints := make([]fcm.Interface, len(scarves))
		for i, s := range scarves {
			fcmPoints[i] = FcmPoint(Point{
				float64(s.Id),
				s.Price,
			})
		}
		priceCluster = FindFCM(fcmPoints, value)
		filtersCount++
	}


	if widthOption, err := strconv.Atoi(c.Query("width")); err == nil && widthOption > 1 {
		value := "min"
		switch widthOption {
		case 3:
			value = "avg"
		case 4:
			value = "max"
		}

		fcmPoints := make([]fcm.Interface, len(scarves))
		for i, s := range scarves {
			fcmPoints[i] = FcmPoint(Point{
				float64(s.Id),
				float64(s.Width),
			})
		}
		widthCluster = FindFCM(fcmPoints, value)
		filtersCount++
	}

	if lengthOption, err := strconv.Atoi(c.Query("length")); err == nil && lengthOption > 1 {
		value := "min"
		switch lengthOption {
		case 3:
			value = "avg"
		case 4:
			value = "max"
		}

		fcmPoints := make([]fcm.Interface, len(scarves))
		for i, s := range scarves {
			fcmPoints[i] = FcmPoint(Point{
				float64(s.Id),
				float64(s.Length),
			})
		}
		lengthCluster = FindFCM(fcmPoints, value)
		filtersCount++
	}

	for _, s := range scarves {
		result := Result{
			Scarf: s,
			TotalWeight: 1,
		}

		if v, ok := priceCluster[s.Id]; ok {
			result.TotalWeight *= v
		}

		if v, ok := lengthCluster[s.Id]; ok {
			result.TotalWeight *= v
		}

		if v, ok := widthCluster[s.Id]; ok {
			result.TotalWeight *= v
		}

		if filtersCount == 1 && result.TotalWeight > 0.5 {
			results = append(results, result)
		} else if filtersCount == 2 && result.TotalWeight > 0.25 {
			results = append(results, result)
		} else if filtersCount == 3 && result.TotalWeight > 0.1 {
			results = append(results, result)
		} else if filtersCount == 0 {
			results = append(results, result)
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i].TotalWeight > results[j].TotalWeight})

	// empty cache
	cache = make([]Scarf, 0)
	c.HTML(http.StatusOK, "final", map[string]interface{}{
		"Rows": results,
	})

}