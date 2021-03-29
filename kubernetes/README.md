# The commands used to deploy the kubedemo application

This demo is for the Vandebron demo, it will only work on our VPN and it comes with the example on how to setup your kubeconfig. 

## External Service setup

```bash
kubectl apply -f ./kubernetes/external/external-deployment.yml
kubectl get svc -n test
kubectl apply -f ./kubernetes/external/external-service.yml
kubectl get svc -n test
kubectl apply -f ./kubernetes/external/external-ingressroute-http.yml
kubectl apply -f ./kubernetes/external/external-ingressroute-https.yml
kubectl -n test get ingressroute
kubectl describe -n test ingressroute kubedemo-ingress-http
curl http://kubedemo.play-backend.zonnecollectief.nl/hello
curl https://kubedemo.play-backend.zonnecollectief.nl/hello
```
