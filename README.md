Fire starter
===

This project makes easier to ADD or REMOVE the access of specifics IPs FROM specifics machines on AWS, via a REST API.


Push to Dockerhub
---

```
make build tag=1
make push tag=1

# Kubernetes service and deployments (changes are needed here)
# TODO - Helm Charts
kubectl create -f kube/*
```

Run
---

You need REDIS with Password enabled and an IAM key with EC2 permissions:

docker run -p 8085:8085 awstools:1 serve -r localhost:6379 -p redis


You can access the REST api as follows:

REST API
---

All the machines are find by tag Names, and by default it enables port 22.

To add the IP for machine called api-1:

```
$ curl -v -XADD localhost:8085/api/v1/add -H "Content-Type: application/json" -d'{"ip": "189.90.98.9", "machine": "api-1"}'
```

To remove the IP address ACL from the machine api-1:

```
$ curl -v -XDELETE localhost:8085/api/v1/ -H "Content-Type: application/json" -d'{"ip": "189.90.98.9", "machine": "api-1"}'
```

Redis Cache
---

After Adding a new IP, it will have a TTL before being removed from the list. It is achieved by Redis PubSub system on events channels.
