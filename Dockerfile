# Build stage
FROM golang:1.19-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /ovh-dynhost-updater
# Final image
FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=build /ovh-dynhost-updater /usr/local/bin/ovh-dynhost-updater
USER 65534:65534
ENTRYPOINT ["ovh-dynhost-updater"]