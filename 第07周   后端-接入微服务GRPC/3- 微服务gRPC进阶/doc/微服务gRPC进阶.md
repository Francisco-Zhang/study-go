## 1、GRPC服务器及客户端

proto带方法生成命令：protoc -I=. --go_out=plugins=grpc,paths=source_relative:gen/go trip.proto
