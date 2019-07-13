package main

import (
	"flag"
	"log"
	"path/filepath"

	versionedclient "github.coaspenmesh/istio-client-go/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	namespace := "huyf"
	restConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create k8s rest client: %s", err)
	}

	ic, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}
	// Test VirtualServices
	vsList, err := ic.NetworkingV1alpha3().VirtualServices(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get VirtualService in %s namespace: %s", namespace, err)
	}
	for i := range vsList.Items {
		vs := vsList.Items[i]
		log.Printf("Index: %d VirtualService Hosts: %+v\n", i, vs.Spec.GetHosts())
	}

	// Test DestinationRules
	drList, err := ic.NetworkingV1alpha3().DestinationRules(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get DestinationRule in %s namespace: %s", namespace, err)
	}
	for i := range drList.Items {
		dr := drList.Items[i]
		log.Printf("Index: %d DestinationRule Host: %+v\n", i, dr.Spec.GetHost())
	}

	// Test Policies
	pList, err := ic.AuthenticationV1alpha1().Policies(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get Policy in %s namespace: %s", namespace, err)
	}
	for i := range pList.Items {
		p := pList.Items[i]
		log.Printf("Index: %d Policy Targets: %+v\n", i, p.Spec.GetTargets())
	}

	// Test MeshPolicies
	mpList, err := ic.AuthenticationV1alpha1().MeshPolicies().List(metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to list MeshPolicies", err)
	}
	for i := range mpList.Items {
		mp := mpList.Items[i]
		log.Printf("Index: %d MeshPolicy Name: %+v\n", i, mp.ObjectMeta.Name)

		// Known broken without the custom marshal/unmarshal code
		log.Printf("Index %d MeshPolicy Value: %+v\n", i, mp.Spec.Policy.Peers)
		_, err := ic.AuthenticationV1alpha1().MeshPolicies().Get(mp.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			log.Fatalf("Failed to get MeshPolicy named %s", mp.ObjectMeta.Name)
		}
	}

	// Test Gateway
	gwList, err := ic.NetworkingV1alpha3().Gateways(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get Gateway in %s namespace: %s", namespace, err)
	}
	for i := range gwList.Items {
		gw := gwList.Items[i]
		for _, s := range gw.Spec.GetServers() {
			log.Printf("Index: %d Gateway servers: %+v\n", i, s)
		}
	}

	// Test ServiceEntry
	seList, err := ic.NetworkingV1alpha3().ServiceEntries(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get ServiceEntry in %s namespace: %s", namespace, err)
	}
	for i := range seList.Items {
		se := seList.Items[i]
		for _, h := range se.Spec.GetHosts() {
			log.Printf("Index: %d ServiceEntry hosts: %+v\n", i, h)
		}
	}
}