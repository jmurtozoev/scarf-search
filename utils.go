package main

var (
	priceOptions = []struct{
		Id int
		Name string
	}{
		{
			Id : 1,
			Name: "",
		},
		{
			Id : 2,
			Name: "Не дорогой",
		},
		{
			Id : 3,
			Name: "Средней цены",
		},
		{
			Id : 4,
			Name: "Дорогой",
		},
	}

	lengthOptions = []struct{
		Id int
		Name string
	}{
		{
			Id : 1,
			Name: "",
		},
		{
			Id : 2,
			Name: "Короткий",
		},
		{
			Id : 3,
			Name: "Средний",
		},
		{
			Id : 4,
			Name: "Длинный",
		},
	}
	widthOptions = []struct{
		Id int
		Name string
	}{
		{
			Id : 1,
			Name: "",
		},
		{
			Id : 2,
			Name: "Узкий",
		},
		{
			Id : 3,
			Name: "Средней",
		},
		{
			Id : 4,
			Name: "Широкий",
		},
	}

	colorOptions = []struct{
		Id int
		Name string
	}{
		{
			Id : 1,
			Name: "",
		},
		{
			Id : 2,
			Name: "Однотонный",
		},
		{
			Id : 3,
			Name: "Средний",
		},
		{
			Id : 4,
			Name: "Цветастый",
		},
	}
)
