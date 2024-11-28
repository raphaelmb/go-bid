FROM golang:1.23 AS build

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o go-bid ./cmd/api/

FROM scratch

WORKDIR /app

COPY --from=build /app/go-bid ./

EXPOSE 3080 

CMD [ "./go-bid" ]
