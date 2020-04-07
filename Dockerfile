FROM golang:buster as builder

  COPY . /logperf/src

  WORKDIR /logperf/src
  RUN make config && \
      make build && \
      mkdir -p /app/web && \
      cp /logperf/src/logperf/logperf /app/logperf && \
      chmod +x /app/logperf && \
      cp -vr /logperf/src/web/* /app/web/


#FROM golang:buster
FROM gcr.io/distroless/cc-debian10:fd0d99e8c54d7d7b2f3dd29f5093d030d192cbbc

    LABEL maintainer="Zack Wine <zwine@synamedia.com>"
    LABEL Description="Logperf docker image" Version="1.1"
    COPY --from=builder /app /app

    EXPOSE 8080
    WORKDIR /app
    CMD ["/app/logperf", "-http"]
