## Ports services

Create ports.json file with ports data or point to existing file in `docker-compose.yml` file and exec:
```
docker-compose up
```

To get port data:
```
curl http://localhost/ports/PORTID
```

To run tests:
```
go test ./...
```