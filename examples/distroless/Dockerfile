FROM golang:1.22 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go get
RUN cd ./examples/distroless && go get && go build -o /out/distroless

FROM scratch AS bin
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/distroless /app/
ENV OTEL_EXPORTER_OTLP_HEADERS="x-honeycomb-team=bogus_key"
EXPOSE 8090
ENTRYPOINT ["/app/distroless"]
