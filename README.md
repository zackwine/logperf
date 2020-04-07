# logperf
Logging performance utility

## Building

Run this once to ensure your go.mod points to your local copy of logperf:

```
make config
```

Build for host:

```
make build
```

Build for linux:

```
make build-linux
```

Build and run

```
make run
# Or
make run-cue
```


## Docker image

### Building the docker image

```
docker build . -t logperf:test
```

### Verifying the docker image HTTP interface is working

```
docker run -p 8080:8080 logperf:test

# From another shell

curl -XPOST -d @test.json 127.0.0.1:8080/v1/api/logperf
```
