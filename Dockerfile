# Start from a Golang v1.16 base image
FROM golang:1.16

# Set the working directory to /app
WORKDIR /app

# Copy the necessary files and directories to the container
COPY . .
# Install the necessary dependencies
# RUN go get github.com/gorilla/mux
# RUN go get github.com/lib/pq

# Build the Go application
RUN go build -o main main.go

# Expose port 8080
EXPOSE 4000

# Set the environment variables
ENV TWILIO_ACCOUNT_SID="ACac72982e1effc9e8656213e939dbd55a"
ENV TWILIO_AUTH_TOKEN="ed24d858ab916a72d5fbd630aa075c3b"
ENV RECEIVER_PHONE="+8801979415290"
ENV SENDER_PHONE="+14155238886"

# Start the Go application
CMD ["./main"]