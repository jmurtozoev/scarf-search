package main

type Scarf struct {
	Id           int     `json:"id"`
	Material     string  `json:"material"`
	Price        float64 `json:"price"`
	Manufacturer string  `json:"manufacturer"`
	Colour       string  `json:"colour"`
	Width        int     `json:"width"`
	Length       int     `json:"length"`
}

type Filter struct {
	Name  string
	Value string
}

type Result struct {
	Scarf       Scarf
	TotalWeight float64
}
