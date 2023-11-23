package main

import (
	"Module/API/initializers"
	"Module/API/view"
)

func init() {
	initializers.LoadEnvVariables()
}
func main() {
	view.Response()
}
