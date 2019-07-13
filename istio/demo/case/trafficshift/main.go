package main

import (
	"log"
	"time"

	vcluster "../common/vcluster"
)

var clus *vcluster.Cluster

func init() {

	clus = vcluster.New()
	clus.Offline()
	clus.Online()
}
func main() {
	shiftTraffic("./config/reviews50v1.yaml")
	shiftTraffic("./config/reviews100v3.yaml")
}

func shiftTraffic(config string) {
	time.Sleep(time.Duration(60) * time.Second)
	toShift := vcluster.LoadSrvYaml(config)
	toShift.ObjectMeta.ResourceVersion = clus.Servics[toShift.GetName()].ObjectMeta.ResourceVersion
	cr, err := clus.IstioCtl.NetworkingV1alpha3().VirtualServices(clus.NameSpace).Update(&toShift)
	clus.Servics[toShift.GetName()].ObjectMeta.ResourceVersion = cr.ObjectMeta.ResourceVersion
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("Service: %+v traffic shifted\n ", cr.Spec.Hosts)
}
