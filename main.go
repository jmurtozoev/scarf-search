package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lhhong/go-fcm/fcm"
	"log"
	"net/http"
	"sort"
)

//var homepageTpl *template.Template
//
//func init() {
//	homepageHTML := assets.MustAssetString("templates/index.html")
//	homepageTpl = template.Must(template.New("homepage_view").Parse(homepageHTML))
//
//}

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html")

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", "")
	})

	router.POST("/search", filterScarves)

	log.Fatal(router.Run(":8080"))

}

type request struct {
	Scarves []Scarf `json:"scarves"`
	Filters []Filter `json:"filters"`
}

func filterScarves(c *gin.Context) {
	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("binding request error: ", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": map[string]interface{}{
				"type": "validation",
				"message": err.Error(),
			},
		})
		return
	}

	results := make([]Result, 0)
	var priceCluster, widthCluster, lengthCluster map[int]float64
	for _, f := range req.Filters {
		// Map your own data slice to the structure implementing fcm.Interface
		if f.Name == "price" {
			fcmPoints := make([]fcm.Interface, len(req.Scarves))
			for i, s := range req.Scarves {
				fcmPoints[i] = FcmPoint(Point{
					float64(s.Id),
					s.Price,
				})
			}
			priceCluster = FindFCM(fcmPoints, f.Value)
		}

		if f.Name == "width" {
			fcmPoints := make([]fcm.Interface, len(req.Scarves))
			for i, s := range req.Scarves {
				fcmPoints[i] = FcmPoint(Point{
					float64(s.Id),
					float64(s.Width),
				})
			}
			widthCluster = FindFCM(fcmPoints, f.Value)
		}

		if f.Name == "length" {
			fcmPoints := make([]fcm.Interface, len(req.Scarves))
			for i, s := range req.Scarves {
				fcmPoints[i] = FcmPoint(Point{
					float64(s.Id),
					float64(s.Length),
				})
			}
			lengthCluster = FindFCM(fcmPoints, f.Value)
		}
	}

	for _, s := range req.Scarves {
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

		if len(req.Filters) == 1 && result.TotalWeight > 0.5 {
			results = append(results, result)
		} else if len(req.Filters) == 2 && result.TotalWeight > 0.25 {
			results = append(results, result)
		} else if len(req.Filters) == 3 && result.TotalWeight > 0.1 {
			results = append(results, result)
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i].TotalWeight > results[j].TotalWeight})

	updatedScarves := make([]Scarf, 0)
	for _, r := range results {
		updatedScarves = append(updatedScarves, r.Scarf)
	}

	c.JSON(http.StatusOK, gin.H{"scarves": updatedScarves})
}