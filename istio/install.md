#1. config istio installtion
```
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.2.2 sh -
cd istio-1.2.2
export PATH=$PWD/bin:$PATH
```
#2. install helm
```
nohup wget https://get.helm.sh/helm-v2.14.1-linux-amd64.tar.gz &
tar -zxvf helm-v2.14.1-linux-amd64.tar.gz
cd linux-amd64/
cp linux-amd64/helm /usr/local/bin
```
#3. install istio
```
kubectl create namespace istio-system
cd istio-1.2.2/
# this step might take a few seconds for the CRDs to be committed in the Kubernetes API-serve
helm template install/kubernetes/helm/istio-init --name istio-init --namespace istio-system | kubectl apply -f -
kubectl get crds | grep 'istio.io\|certmanager.k8s.io' | wc -lv
helm template install/kubernetes/helm/istio --name istio --namespace istio-system | kubectl apply -f -
```


