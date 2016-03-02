package main

import (
	"TodoBackend/controllers"
	_ "TodoBackend/utils"
)

/*
Simplifying assumptions (for now):
- users only have one device
- no permissions
- no push

Things to figure out:
- synchronization
*/
func main() {
	r := controllers.RegisterRoutes()
	r.Run(":8080")
}
