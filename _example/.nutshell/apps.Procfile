user: go run apps/main.go --app=user --http_port=8001 --grpc_port=9001
rel: go run apps/main.go --app=rel --grpc_port=9002 --user_grpc_addr=127.0.0.1:9001
im: go run apps/main.go --app=im --http_port=8002 --user_grpc_addr=127.0.0.1:9001 --rel_grpc_addr=127.0.0.1:9002
