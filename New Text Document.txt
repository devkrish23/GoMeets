# Use the official Golang image with the specified version
FROM golang:1.19-alpine

# Set the working directory inside the container
WORKDIR /src

# Copy the Go module files (go.mod and go.sum) to the working directory
COPY go.mod go.sum ./
# Download and cache Go dependencies
RUN go mod download
# Copy the rest of the application source code into the working directory
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/app ./cmd

FROM alpine
WORKDIR /src

# Copy the views and assets directories from the previous stage to the current stage

COPY --from=0 /bin/app /bin/app
COPY --from=0 /src/views /src/views
COPY --from=0 /src/assets /src/assets

ENTRYPOINT ["/bin/app"]