FROM node:lts-alpine as app

WORKDIR /app

ARG VITE_SERVER_URL="http://localhost:8090"

ENV VITE_SERVER_URL $VITE_SERVER_URL

ADD /app ./

RUN npm i

RUN npm run build

# --

FROM golang:1.18-alpine as builder

RUN apk add git

WORKDIR /app

COPY . .

COPY --from=app /app/dist ./app/dist

COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -a -installsuffix cgo -o main .

# --

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

COPY --from=builder /app/main ./

VOLUME /app/pb_data

EXPOSE 8090

CMD ["./main", "serve", "--http=0.0.0.0:8090"]