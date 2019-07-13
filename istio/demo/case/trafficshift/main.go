package main

import (
	"time"

	extApi "../common/api"
)

func init() {
	extApi.DefaultServiceCollection.Online()

}
func main() {
	reviews50v1 := extApi.VirtualServiceFromYaml("./config/reviews50v1.yaml")
	reviews100v3 := extApi.VirtualServiceFromYaml("./config/reviews100v3.yaml")
	time.Sleep(time.Duration(60) * time.Second)
	reviews50v1.Online()
	time.Sleep(time.Duration(60) * time.Second)
	reviews100v3.Online()
}
