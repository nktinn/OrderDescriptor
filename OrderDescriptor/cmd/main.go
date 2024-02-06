package main

import "github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/app"

const configPath = "OrderDescriptor/config/config.yaml"

func main() {
	app.Run(configPath)
}
