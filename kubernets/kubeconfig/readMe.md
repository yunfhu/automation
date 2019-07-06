#1. create admin service acoount and bind with cluster role
```
apiVersion: v1
kind: ServiceAccount
metadata:
  name: huyf
  namespace: kube-system

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: yugy
  namespace: kube-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: hull
  namespace: kube-system
---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding 
metadata: 
  name: admin-user
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: huyf
  namespace: kube-system
- kind: ServiceAccount
  name: yugy
  namespace: kube-system
- kind: ServiceAccount
  name: hull
  namespace: kube-system


```

# 2. config.yaml
```
apiVersion: v1
kind: Config
preferences: {}

clusters:
- cluster:
  name: AsiaInfo

users:
- name: huyf

contexts:
- context:
  name: AsiaInfo
```
#3. 设置config.yaml

kubectl config --kubeconfig=./config.yaml  set-cluster AsiaInfo --server=https://10.21.20.88:6443 --insecure-skip-tls-verify
kubectl config  --kubeconfig=./config.yaml set-credentials huyf --token=$(kubectl get secret huyf-token-n8w5k -n kube-system -o jsonpath={.data.token} | base64 -D)
kubectl config --kubeconfig=./config.yaml set-context AsiaInfo --cluster=AsiaInfo --namespace=default --user=huyf
kubectl config use-context AsiaInfo










