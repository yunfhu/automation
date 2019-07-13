package api

import (
	"flag"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"runtime"

	"github.coaspenmesh/istio-client-go/pkg/client/clientset/versioned"
	versionedclient "github.coaspenmesh/istio-client-go/pkg/client/clientset/versioned"
	"github.com/aspenmesh/istio-client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/yaml"
)

var (
	nameSpace                = "huyf"
	IstioCtl                 *versioned.Clientset
	DefaultServiceCollection VirtualServiceCollection
)

type lifeCycle interface {
	Online()
	Offline()
}
type VirtualService struct {
	MetaService *v1alpha3.VirtualService
}
type DestinationRule struct {
	MetaRule *v1alpha3.DestinationRule
}
type VirtualServiceCollection struct {
	Servics (map[string]*VirtualService)
}

func (v *VirtualService) Online() {
	metaService := v.MetaService
	var cs *v1alpha3.VirtualService
	var err error
	ps, geterr := IstioCtl.NetworkingV1alpha3().VirtualServices(nameSpace).Get(metaService.GetObjectMeta().GetName(), v1.GetOptions{})
	if geterr != nil {
		cs, err = IstioCtl.NetworkingV1alpha3().VirtualServices(nameSpace).Create(metaService)
	} else {
		metaService.ObjectMeta.SetResourceVersion(ps.GetObjectMeta().GetResourceVersion())
		cs, err = IstioCtl.NetworkingV1alpha3().VirtualServices(nameSpace).Update(metaService)
	}
	CheckErr(err)
	log.Printf("Service: %+v  online\n ", cs.Spec.Hosts)
}

func (v *VirtualService) Offline() {
	metaService := v.MetaService
	var cs *v1alpha3.VirtualService
	var err error
	ps, geterr := IstioCtl.NetworkingV1alpha3().VirtualServices(nameSpace).Get(metaService.GetObjectMeta().GetName(), v1.GetOptions{})
	if geterr != nil {
		cs, err = IstioCtl.NetworkingV1alpha3().VirtualServices(nameSpace).Create(metaService)
	} else {
		metaService.ObjectMeta.SetResourceVersion(ps.GetObjectMeta().GetResourceVersion())
		cs, err = IstioCtl.NetworkingV1alpha3().VirtualServices(nameSpace).Update(metaService)
	}
	CheckErr(err)
	log.Printf("Service: %+v  online\n ", cs.Spec.Hosts)
}

func (cls *VirtualServiceCollection) Online() {
	for _, v := range cls.Servics {
		v.Online()
	}
}
func (cls *VirtualServiceCollection) Offline() {
	for _, v := range cls.Servics {
		v.Offline()
	}
}
func (rule *DestinationRule) Online() {
	metdaRule := rule.MetaRule
	var cr *v1alpha3.DestinationRule
	var err error
	pr, err := IstioCtl.NetworkingV1alpha3().DestinationRules(nameSpace).Get(metdaRule.GetName(), v1.GetOptions{})
	if err != nil {
		cr, err = IstioCtl.NetworkingV1alpha3().DestinationRules(nameSpace).Create(metdaRule)
	} else {
		metdaRule.ObjectMeta.ResourceVersion = pr.ObjectMeta.ResourceVersion
		cr, err = IstioCtl.NetworkingV1alpha3().DestinationRules(nameSpace).Update(metdaRule)
	}
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("rule :%s online \n ", cr.GetName())

}

func (rule *DestinationRule) Offline() {
	metaRule := rule.MetaRule
	err := IstioCtl.NetworkingV1alpha3().DestinationRules(nameSpace).Delete(metaRule.GetName(), &v1.DeleteOptions{})
	CheckErr(err)
	log.Printf("rule: %s offline", metaRule.GetName())
}

func init() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	restConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create k8s rest client: %s", err)
	}

	ctl, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}
	IstioCtl = ctl
	_, filename, _, _ := runtime.Caller(0)
	DefaultServiceCollection = VirtualServiceCollectionFromFolder(path.Join(path.Dir(filename), "./config/"))

}

func CheckErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func DestinationRuleFromYaml(path string) DestinationRule {
	var rule = v1alpha3.DestinationRule{}
	buffer, err := ioutil.ReadFile(path)
	CheckErr(err)
	j2, err := yaml.YAMLToJSON(buffer)
	CheckErr(err)
	err = yaml.Unmarshal(j2, &rule)
	CheckErr(err)
	return DestinationRule{
		MetaRule: &rule,
	}
}

func VirtualServiceCollectionFromFolder(folder string) VirtualServiceCollection {
	var svs = make(map[string]*VirtualService)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Panicf("folder:%s read failed", folder)
	}
	for _, f := range files {
		vs := VirtualServiceFromYaml(path.Join(folder, f.Name()))
		svs[vs.MetaService.GetName()] = &vs
	}
	return VirtualServiceCollection{Servics: svs}

}

func VirtualServiceFromYaml(path string) VirtualService {
	var v = v1alpha3.VirtualService{}
	buffer, err := ioutil.ReadFile(path)
	CheckErr(err)
	j2, err := yaml.YAMLToJSON(buffer)
	CheckErr(err)
	err = yaml.Unmarshal(j2, &v)
	CheckErr(err)
	return VirtualService{
		MetaService: &v,
	}
}
