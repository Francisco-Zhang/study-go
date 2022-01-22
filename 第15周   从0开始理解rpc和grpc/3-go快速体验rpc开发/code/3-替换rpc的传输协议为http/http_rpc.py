import requests

request = {
    "id":0,
    "params":["imooc"],
    "method": "HelloService.Hello"
}

rsp= requests.post("http://localhost:1234/jsonrpc",json=request)
print(rsp.text)