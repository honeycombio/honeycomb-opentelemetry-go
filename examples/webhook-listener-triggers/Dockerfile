FROM golang:1.19-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go get
RUN go build -o /out/webhook-listener-triggers .

FROM scratch AS bin
WORKDIR /app
COPY --from=build /out/webhook-listener-triggers /app/
ENV HONEYCOMB_API_KEY=bogus_key
EXPOSE 8090
ENTRYPOINT ["/app/webhook-listener-triggers"]
