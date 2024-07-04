# How to deploy Golang application on Kubernetes 

## Building the Go application 
- Using Dockerfile
- Using Ko
```
KO_DOCKER_REPO=ttl.sh ko publish .
```


Using Buildsafe
## Installing nginx ingress controller

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml
```

## Installing Cert Manager
```
  kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.1/cert-manager.yaml
```

## Create Deployment and Service 
```
cat << EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-server-deployment
  labels:
    app: go-server
spec:
  replicas: 2
  selector:
    matchLabels:
      app: go-server
  template:
    metadata:
      labels:
        app: go-server
    spec:
      containers:
        - name: go-server
          image: image
          ports:
            - containerPort: 8080
EOF
```
```
cat << EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: go-server-service
  labels:
    app: go-server
spec:
  type: NodePort
  ports:
    - port: 80
      targetPort: 8080
  selector:
    app: go-server
EOF
```

## Create Ingress
```
cat << EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    kubernetes.io/ingress.class: nginx
  name: demo-app-ingress
spec:
  rules:
  - host: go.kubesimplify.com
    http:
      paths:
      - backend:
          service:
            name: go-server-service
            port:
              number: 80
        path: /
        pathType: Prefix
  tls:
  - hosts:
    - go.kubesimplify.com
    secretName: demo
EOF

```

## Create Certificate and Clusterissuer
```
cat << EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: saiym911@gmail.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: demo
spec:
  secretName: demo
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  commonName: go.kubesimplify.com
  dnsNames:
  - go.kubesimplify.com
EOF
```
