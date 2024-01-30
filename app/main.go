package main

import (
	"your-project/internal/router"
)

func main() {
	r := router.SetupRouter()
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
