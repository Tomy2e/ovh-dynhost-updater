# Build stage
FROM golang:1.17-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /ovh-dynhost-updater
# Final image
FROM alpine:3.14
WORKDIR /
COPY --from=build /ovh-dynhost-updater /usr/local/bin/ovh-dynhost-updater
USER 10001:10001
ENTRYPOINT ["ovh-dynhost-updater"]