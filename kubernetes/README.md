# The commands used to deploy the kubedemo application

This demo is for the Vandebron demo, it will only work on our VPN and it comes with the example on how to setup your kubeconfig. 

## External Service setup

```bash
kubectl get deployment -n test
kubectl apply -f ./kubernetes/external/external-deployment.yml
kubectl get deployment -n test
kubectl get svc -n test
kubectl apply -f ./kubernetes/external/external-service.yml
kubectl get svc -n test
kubectl apply -f ./kubernetes/external/external-ingressroute-http.yml
kubectl apply -f ./kubernetes/external/external-ingressroute-https.yml
kubectl -n test get ingressroute
kubectl describe -n test ingressroute kubedemo-ingress-http
curl -I http://kubedemo.play-backend.zonnecollectief.nl/hello
curl https://kubedemo.play-backend.zonnecollectief.nl/hello
```

## Internal Service url

```bash
kubectl run --namespace test tester-$RANDOM --rm --tty -i  --image anldisr/curl:latest -- curl kubedemo.test.svc.cluster.local:8090/hello
```

## Clean up

```bash
kubectl delete -f ./kubernetes/external/external-deployment.yml
kubectl delete -f ./kubernetes/external/external-service.yml
kubectl delete -f ./kubernetes/external/external-ingressroute-http.yml
kubectl delete -f ./kubernetes/external/external-ingressroute-https.yml
```
