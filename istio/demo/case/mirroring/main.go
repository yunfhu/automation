package main

import (
	"time"

	extApi "../common/api"
)

func init() {
	extApi.DefaultServiceCollection.Online()

}
func main() {
	reviews100v1 := extApi.VirtualServiceFromYaml("./config/reviews100v1.yaml")
	reviewsv1v2Rule := extApi.DestinationRuleFromYaml("./config/reviewsv1v2Rule.yaml")
	reviewsv2mirrorv1 := extApi.VirtualServiceFromYaml("./config/reviewsv2mirrorv1.yaml")
	reviews100v1.Online()
	reviewsv1v2Rule.Online()
	time.Sleep(time.Duration(15) * time.Second)
	reviewsv2mirrorv1.Online()
}
