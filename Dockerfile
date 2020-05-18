FROM golang:latest

# Copy the local package files to the container's workspace.
RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go get github.com/gin-gonic/gin
RUN go build -o gowebapp .

# Run the outyet command by default when the container starts.
ENTRYPOINT /app/gowebapp

# Document that the service listens on port 8081.
EXPOSE 8080

