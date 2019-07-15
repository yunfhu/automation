package main

import (
	"time"

	extApi "../common/api"
)

func init() {
	extApi.DefaultServiceCollection.Online()
}
func main() {
	circuitBreak()
}
func circuitBreak() {
	productpagebin := extApi.DestinationRuleFromYaml("./productpagebin.yaml")
	productpagebin.Online()
	time.Sleep(time.Duration(30) * time.Second)
	defer productpagebin.Offline()
}
