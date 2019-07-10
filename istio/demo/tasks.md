#1. Request Routing
The Istio Bookinfo sample consists of four separate microservices, each with multiple versions. Three different versions of one of the microservices, reviews, have been deployed and are running concurrently. To illustrate the problem this causes, access the Bookinfo app’s /productpage in a browser and refresh several times. You’ll notice that sometimes the book review output contains star ratings and other times it does not. This is because without an explicit default service version to route to, Istio routes requests to all available versions in a round robin fashion.
The initial goal of this task is to apply rules that route all traffic to v1 (version 1) of the microservices. Later, you will apply a rule to route traffic based on the value of an HTTP request header
1. refresh the [bookInfo page](http://10.21.19.77:31380/productpage)
2. apply the virtual services
```
kubectl apply -f istio/demo/samples/bookinfo/networking/virtual-service-all-v1.yaml
```
3. Display the defined routes
```
kubectl get virtualservices -o yaml
```
```
kubectl get destinationrules -o yaml
```
4. uninstall the virtual services
```
kubectl delete -f istio/demo/samples/bookinfo/networking/virtual-service-all-v1.yaml
```
#2.Fault Injection
1. Apply application version routing 
```
kubectl apply -f istio/demo/samples/bookinfo/networking/virtual-service-all-v1.yaml
kubectl apply -f istio/demo/samples/bookinfo/networking/virtual-service-reviews-test-v2.yaml
```
With the above configuration, this is how requests flow:
productpage → reviews:v2 → ratings (only for user jason)
productpage → reviews:v1 (for everyone else)
2. Create a delay fault injection rule to delay traffic coming from the test user jason.
```
kubectl apply -f istio/demo/samples/bookinfo/networking/virtual-service-ratings-test-delay.yaml
```
3.  Create a abort fault injection rule to delay traffic coming from the test user jason.
```
kubectl apply -f istio/demo/samples/bookinfo/networking/virtual-service-ratings-test-abort.yaml
```





