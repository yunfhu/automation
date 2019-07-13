package main

import (
	"io/ioutil"
	"log"
	"time"

	util "../common/utils"
	vcluster "../common/vcluster"
	v1alpha3 "github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var clus *vcluster.Cluster

func init() {

	clus = vcluster.New()
	clus.Offline()
	clus.Online()
}
func main() {
	circuitBreak("./productpagebin.yaml")
}
func circuitBreak(path string) {
	rule := loadRuleYaml(path)
	onlineRule(&rule)
	time.Sleep(time.Duration(30) * time.Second)
	defer offlineRule(&rule)
}

func onlineRule(rl *v1alpha3.DestinationRule) {
	var cr *v1alpha3.DestinationRule
	var err error
	pr, err := clus.IstioCtl.NetworkingV1alpha3().DestinationRules(clus.NameSpace).Get(rl.GetName(), v1.GetOptions{})
	if err != nil {
		cr, err = clus.IstioCtl.NetworkingV1alpha3().DestinationRules(clus.NameSpace).Create(rl)
	} else {
		rl.ObjectMeta.ResourceVersion = pr.ObjectMeta.ResourceVersion
		cr, err = clus.IstioCtl.NetworkingV1alpha3().DestinationRules(clus.NameSpace).Update(rl)
	}
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("rule :%s online \n ", cr.GetName())

}

func offlineRule(rl *v1alpha3.DestinationRule) {
	err := clus.IstioCtl.NetworkingV1alpha3().DestinationRules(clus.NameSpace).Delete(rl.GetName(), &v1.DeleteOptions{})
	util.Checkerr(err)
	log.Printf("rule: %s offline", rl.GetName())
}

func loadRuleYaml(path string) v1alpha3.DestinationRule {
	var rule = v1alpha3.DestinationRule{}
	buffer, err := ioutil.ReadFile(path)
	util.Checkerr(err)
	j2, err := yaml.YAMLToJSON(buffer)
	util.Checkerr(err)
	err = yaml.Unmarshal(j2, &rule)
	util.Checkerr(err)
	return rule
}
