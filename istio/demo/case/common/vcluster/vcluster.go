package cls

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
	"github.com/ghodss/yaml"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var yamlfolder = "./config/"

var Clus *Cluster

type lifeCycle interface {
	Online()
	Offline()
}

type Cluster struct {
	NameSpace string
	IstioCtl  *versioned.Clientset
	Servics   (map[string]*v1alpha3.VirtualService)
}

func (cls *Cluster) Offline() {
	for k := range cls.Servics {
		err := cls.IstioCtl.NetworkingV1alpha3().VirtualServices(cls.NameSpace).Delete(k, &v1.DeleteOptions{})
		if nil != err {
		}
		// checkerr(err)
		log.Printf("Service: %+v  offline\n ", k)
	}
}
func (cls *Cluster) Online() {
	for _, v := range cls.Servics {
		cr, err := cls.IstioCtl.NetworkingV1alpha3().VirtualServices(cls.NameSpace).Create(v)
		checkerr(err)
		v.ObjectMeta.ResourceVersion = cr.ObjectMeta.ResourceVersion
		log.Printf("Service: %+v  online\n ", cr.Spec.Hosts)
	}
}

func New() *Cluster {
	cls := Cluster{
		NameSpace: "huyf",
	}
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
	cls.IstioCtl = ctl
	_, filename, _, _ := runtime.Caller(0)
	log.Printf(path.Join(path.Dir(filename), yamlfolder))
	cls.Servics = BatchLoadSrvYaml(path.Join(path.Dir(filename), yamlfolder))
	return &cls
}

func BatchLoadSrvYaml(folder string) map[string]*v1alpha3.VirtualService {
	var svs = make(map[string]*v1alpha3.VirtualService)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Panicf("folder:%s read failed", folder)
	}
	for _, f := range files {
		vs := LoadSrvYaml(path.Join(folder, f.Name()))
		svs[vs.GetName()] = &vs
	}
	return svs

}

func LoadSrvYaml(path string) v1alpha3.VirtualService {
	var v = v1alpha3.VirtualService{}
	buffer, err := ioutil.ReadFile(path)
	checkerr(err)
	j2, err := yaml.YAMLToJSON(buffer)
	checkerr(err)
	err = yaml.Unmarshal(j2, &v)
	checkerr(err)
	return v
}

func checkerr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}
