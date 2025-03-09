# Elastic Local

### Run elastic server
```shell
docker run --name es01-test --net elastic -p 127.0.0.1:9200:9200 -p 127.0.0.1:9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:8.5.3
```

### Run Kibana
```shell
docker run --name kib01-test --net elastic -p 127.0.0.1:5601:5601 -e "ELASTICSEARCH_HOSTS=http://es01-test:9200" docker.elastic.co/kibana/kibana:7.17.3
```

### Download And Run apm server
```shell
# download
curl -L -O https://artifacts.elastic.co/downloads/apm-server/apm-server-7.17.8-darwin-x86_64.tar.gz
tar xzvf apm-server-7.17.8-darwin-x86_64.tar.gz
cd apm-server-7.17.8-darwin-x86_64/

#run
./apm-server -e
```