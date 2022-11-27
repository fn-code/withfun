docker build -t envoy:http .

docker run -p 8080:8080 envoy:http