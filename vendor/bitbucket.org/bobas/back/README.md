# Development

## Running integration tests

Create docker machine:
```
docker-machine create --driver virtualbox nildev
```

Setup docker environment:
```
docker-machine start nildev
eval $(docker-machine env nildev)
```

Run required containers:
```
docker-compose -f docker-compose-dev.yml up -d
```

Execute tests

Provisioning is happening inside `TestMain`
```
ND_IP_PORT=$(dm ip nildev):27017 go test -v -tags integration
```

# Running server

```
ND_DATABASE_NAME=back ND_MONGODB_URL=mongodb://192.168.99.100:27017/back ND_STORAGE=mongodb ND_PASS="xxx" ./run
```

## Project Details

### Release Notes

See the [releases tab](https://github.com/nildev/back/releases) for more information on each release.

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and contacting developers via IRC and mailing lists.

### License

Project is released under the MIT license. See the [LICENSE](LICENSE) file for details.

Specific components of project use code derivative from software distributed under other licenses; in those cases the appropriate licenses are stipulated alongside the code.