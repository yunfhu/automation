package main

// need open page http://10.21.19.77:31380/productpage show the affection,
//pay attention to the "Book Reviews" block,
//this case conbined  fault injection function with request timeout function

import (
	"time"

	extApi "../common/api"
)

func init() {
	extApi.DefaultServiceCollection.Online()

}
func main() {
	//direct request to reviews service version 2
	reviewsv2 := extApi.VirtualServiceFromYaml("./config/reviewsv2.yaml")
	//inject delay 2s to  rating version 1 services
	ratingsv1delay2s := extApi.VirtualServiceFromYaml("./config/ratingsv1delay2s.yaml")
	//set the reviews v2 timeout as 1s
	reviewsv2timeout1s := extApi.VirtualServiceFromYaml("./config/reviewsv2timeout1s.yaml")
	reviewsv2.Online()
	ratingsv1delay2s.Online()
	time.Sleep(time.Duration(30) * time.Second)
	reviewsv2timeout1s.Online()
}
