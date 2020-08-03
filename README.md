# Scanley

Scanley is a port scanner written in go.

# Getting Started

## Build the image
1. `make build`

## Deployment
### Kind/K3D

1. Apply the Manifest
```
$ kubectl apply -f k8s/stack.yml
deployment.apps/scanley created
service/scanley created
ingress.extensions/scanley created
```

2. Open a browser and navigate to http://127.0.0.1.nip.io

*note: If you don't have an ingress controller, you can also tunnel to the application by running the following command: `kubectl port-forward svc/scanley 8080:80` , and navigate to http://localhost:8080*


## Run in Docker
1. Run the image 
    * `docker run -p 8080:8080 sgryczan/scanley:0.0.1`
2. Navigate to http://localhost:8080

# Usage

## Anatomy of a scan request

Example:
```
{
  "hosts": [
    "scanme.nmap.org",
    "google.com",
    "amazon.com"
  ],
  "ports": [
    25,
    80,
    443
  ],
  "port-ranges": [
    "31000-32000"
  ],
  "connect-timeout": 3000,
  "max-concurrent-hosts": 10,
  "threads": 500
}
```

| Field | Required | Description |
| ----------- | ----------- | ----------- |
| hosts | True | One or more host names/IP addresses to scan |
| ports | True | List of ports to scan |
| port-panges | False | List of port ranges to scan. Syntax: **\<start port>-\<end port>** |
| connect-timeout | false | Port connection timeout. If port does not respond within this period, it is marked closed. Increasing this will increase scan accuracy, but scans will take longer to complete. |
| max-concurrent-hosts | false | Max number of hosts to scan concurrently. |
| threads | false | Number of threads per scan, per host |


## Run a basic scan
Request:
```
$ curl -X POST "http://127.0.0.1.nip.io/scan" \
    -H  "accept: text/plain" \
    -H  "Content-Type: application/json" \
    -d '{"hosts": ["scanme.nmap.org"], "ports": [25, 80, 443]}'

```
Respose:
```
{
  "scan-id": "bsjntedvr1d5ia8d3m60"
}
```

## List scans
Request:
```
curl -X GET "http://127.0.0.1.nip.io/scan" -H  "accept: application/json"
```
Response: 
```
{
  "count": 4,
  "scans": [
    "bsjnsitvr1d5ia8d3m5g",
    "bsjntedvr1d5ia8d3m60",
    "bsjnublvr1d5ia8d3m6g",
    "bsjnubtvr1d5ia8d3m70"
  ]
}
```

## Get scan results
Request:
```
$ curl -X GET "http://127.0.0.1.nip.io/scan/bsjnvv5vr1d5ia8d3m7g" -H  "accept: application/json"
```
Response:
```
{
  "date": "2020-08-03T03:12:21Z",
  "results": [
    {
      "host": "scanme.nmap.org",
      "openports": [
        22,
        80,
        9929,
        31337
      ]
    }
  ]
}
```