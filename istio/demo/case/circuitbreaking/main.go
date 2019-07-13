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
	rule := extApi.DestinationRuleFromYaml("./productpagebin.yaml")
	rule.Online()
	time.Sleep(time.Duration(30) * time.Second)
	defer rule.Offline()
}
