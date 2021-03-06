#1.login to the k8s master node
#2. define the credentials you want to use as the Kiali username and passphrase
```
KIALI_USERNAME=$(read -p 'Kiali Username: ' uval && echo -n $uval | base64)
KIALI_PASSPHRASE=$(read -sp 'Kiali Passphrase: ' pval && echo -n $pval | base64)
```
#3.create a secret
```
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: kiali
  namespace: istio-system
  labels:
    app: kiali
type: Opaque
data:
  username: $KIALI_USERNAME
  passphrase: $KIALI_PASSPHRASE
EOF
```
#4.intall via helm
```
cd istio-1.2.2
helm template --set kiali.enabled=true install/kubernetes/helm/istio --name istio --namespace istio-system > $HOME/istio.yaml
kubectl apply -f $HOME/istio.yaml
```
> This task does not discuss Jaeger and Grafana. If you already installed them in your cluster and you want to see how Kiali integrates with them, you must pass additional arguments to the helm command, for example:
```
helm template \
    --set kiali.enabled=true \
    --set "kiali.dashboard.jaegerURL=http://jaeger-query:16686" \
    --set "kiali.dashboard.grafanaURL=http://grafana:3000" \
    install/kubernetes/helm/istio \
    --name istio --namespace istio-system > $HOME/istio.yaml
kubectl apply -f $HOME/istio.yaml
```
> 

#5. To open the Kiali UI, execute the following command in your Kubernetes environment:
```
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=kiali -o jsonpath='{.items[0].metadata.name}') 20001:20001
```
Visit http://localhost:20001/kiali/console in your web browser.




