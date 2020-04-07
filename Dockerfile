FROM golang:buster as builder

  COPY . /logperf/src

  WORKDIR /logperf/src
  RUN make config && \
      make build && \
      mkdir -p /app/web && \
      cp /logperf/src/logperf/logperf /app/logperf && \
      chmod +x /app/logperf && \
      cp -vr /logperf/src/web/* /app/web/




FROM golang:buster

    LABEL maintainer="Zack Wine <zwine@synamedia.com>"
    LABEL Description="Logperf docker image" Version="1.1"
    COPY --from=builder /app /app

    EXPOSE 8080
    WORKDIR /app
    CMD ["/app/logperf", "-http"]