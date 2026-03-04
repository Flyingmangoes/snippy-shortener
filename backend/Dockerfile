# Set docker base image 
# Set root dir
# Install required dependency
# Copy all files to container
# Inject environment variable

FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go.mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Running the app
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/server .
EXPOSE 8080
CMD [ "/server" ]
