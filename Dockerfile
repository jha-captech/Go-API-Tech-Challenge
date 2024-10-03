FROM golang:1.23-alpine AS builder

WORKDIR /gotechchallenge


# Copy source code.
COPY go.mod go.sum ./
COPY cmd /gotechchallenge/cmd
COPY internal /gotechchallenge/internal

# Download dependencies.
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build /gotechchallenge/cmd/api/main.go 

FROM alpine:3.19 AS publish

COPY --from=builder /gotechchallenge/main .

EXPOSE 8080

ENTRYPOINT [ "./main" ]
