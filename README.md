# Go restful composite microservice

## Description
This service makes use of the scatter and gather pattern. It executes http requests in parallel and the merges the reslutant data
It uses the following strategy
- Fail Once : Fail all if only one request fails (default)
- Fail None : Ignore any failure
- Pass Once : Pass if at least one request passes

## Usage 

```bash
# post data schema (local service)
curl -d'{"strategy": "failonce", "requests": [{"method":"GET","url":"http://127.0.0.1:8081/api/v1/service","payload":""},{"method":"POST","url":"http://127.0.0.1:8082/api/v1/service","payload":"{wtf:asshole-toolbox}"},{"method":"POST","headers":"api-key:sadhsjahdsahd,content-type:application/json","url":"http://127.0.0.1:8083/api/v1/service","payload":"{wtf: reedge}"} ]}' http://127.0.0.1:9001/api/v1/composite
```
