FROM golang:buster AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD ./ ./
RUN go build -o /images-reverse-proxy


FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /images-reverse-proxy /images-reverse-proxy
USER nonroot:nonroot
ENTRYPOINT ["/images-reverse-proxy"]
