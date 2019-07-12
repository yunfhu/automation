#1. Determining the ingress IP and port
```
# determine if your Kubernetes cluster is running in an environment that supports external load balancers
kubectl get svc istio-ingressgateway -n istio-system
```
As the the EXTERNAL-IP value is <none> (or perpetually <pending>), environment does not provide an external load balancer for the ingress gateway. In this case, you can access the gateway using the service’s node port.
```
#set INGRESS_PORT
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
#set INGRESS_HOST
export INGRESS_HOST=$(kubectl get po -l istio=ingressgateway -n istio-system -o jsonpath='{.items[0].status.hostIP}')
export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT
```
#2. deploy sample services
##2.1 deploy the app
```
#enable namespace isto injection
kubectl label namespace huyf istio-injection=enabled
#deploy sample microservice
kubectl apply -f istio/demo/samples/bookinfo/platform/kube/bookinfo.yaml
```
## 2.2 To confirm that the Bookinfo application is running, send a request to it by a curl command from some pod, for example from ratings:
```
 kubectl exec -it $(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}') -c ratings -- curl productpage:9080/productpage | grep -o "<title>.*</title>"
 ```
## 2.3. Define the gateway for the application
```
kubectl apply -f istio/demo/samples/bookinfo/networking/bookinfo-gateway.yaml
```
## 2.4 Confirm the app is accessible from outside the cluster 
```
curl -s http://${GATEWAY_URL}/productpage | grep -o "<title>.*</title>"
<title>Simple Bookstore App</title>
```

## 2.5. Apply default destination rules
* If you did not enable mutual TLS, execute this command:
 ```
 kubectl apply -f istio/demo/samples/bookinfo/networking/destination-rule-all.yaml
 ```
* If you did enable mutual TLS, execute this command:
```
kubectl apply -f istio/demo/samples/bookinfo/networking/destination-rule-all-mtls.yaml
```
Wait a few seconds for the destination rules to propagate.You can display the destination rules with the following command:
```
kubectl get destinationrules -o yaml
```

more infomation pls read [bookinfo examples](https://istio.io/docs/examples/bookinfo/)









