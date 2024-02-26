### Building and running your application

When you're ready, start your application by running:
`docker compose up --build`.

Application expose grpc endpoints on 50051 port

### Testing

Run this command to get grpcui image deployed
`docker run -eGRPCUI_SERVER=172.17.0.1:50051 -p8080:8080 wongnai/grpcui`

grpc ui client would be available at: http://127.0.0.1:8080/