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


curl -d'{"strategy": "failonce", "requests": [{"method":"GET","url":"http://192.168.0.126:9000/api/v1/service","payload":""},{"method":"POST","url":"http://192.168.0.127:9000/api/v1/service","payload":"{wtf:stupid-simple}"},{"method":"POST","headers":"api-key:sadhsjahdsahd,content-type:application/json","url":"http://192.168.0.128:9000/api/v1/service","payload":"{wtf: dude}"},{"method":"POST","url":"http://192.168.0.129:9000/api/v1/service","payload":"{wtf:ah-toolbox}"}, {"method":"POST","url":"http://192.168.0.130:9000/api/v1/service","payload":"{wtf:ah-toolbox}"} ]}' http://127.0.0.1:9001/api/v1/composite

# results

{
	"message": "ExecuteComposite data successfully aggregated",
	"statuscode": "200",
	"status": "OK",
	"results": {
		"lastupdate": 1586636568913418086,
		"strategy": "failonce",
		"requests": [
			{
				"method": "GET",
				"url": "http://192.168.0.126:9000/api/v1/service",
				"payload": "{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"\"\n}",
				"headers": "",
				"message": "Request called successfully",
				"statuscode": "200",
				"status": "OK"
			},
			{
				"method": "POST",
				"url": "http://192.168.0.127:9000/api/v1/service",
				"payload": "{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"{wtf:stupid-simple}\"\n}",
				"headers": "",
				"message": "Request called successfully",
				"statuscode": "200",
				"status": "OK"
			},
			{
				"method": "POST",
				"url": "http://192.168.0.128:9000/api/v1/service",
				"payload": "{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"{wtf: dude}\"\n}",
				"headers": "api-key:sadhsjahdsahd,content-type:application/json",
				"message": "Request called successfully",
				"statuscode": "200",
				"status": "OK"
			},
			{
				"method": "POST",
				"url": "http://192.168.0.129:9000/api/v1/service",
				"payload": "{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"{wtf:ah-toolbox}\"\n}",
				"headers": "",
				"message": "Request called successfully",
				"statuscode": "200",
				"status": "OK"
			},
			{
				"method": "POST",
				"url": "http://192.168.0.130:9000/api/v1/service",
				"payload": "{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"{wtf:ah-toolbox}\"\n}",
				"headers": "",
				"message": "Request called successfully",
				"statuscode": "200",
				"status": "OK"
			}
		],
		"mergedcontent": [
			"{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"\"\n}",
			"{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"{wtf:stupid-simple}\"\n}",
			"{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"{wtf: dude}\"\n}",
			"{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"{wtf:ah-toolbox}\"\n}",
			"{\n\t\"status\": \"OK\",\n\t\"statuscode\": \"200\",\n\t\"message\": \"Request Successfulli\",\n\t\"result\": \"{wtf:ah-toolbox}\"\n}"
		]
	}
}
```


