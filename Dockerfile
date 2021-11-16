# Start from golang v1.13 base image
FROM golang:1.14

# Add Maintainer Info
# LABEL maintainer="LVC"

# Set the Current Working Directory inside the container
WORKDIR /go/src/neurone-am-simulator-v2/

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container

COPY go.mod .
ENV GO111MODULE=on


RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o neurone-am-simulator-v2 .


######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash 


WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=0 /go/src/neurone-am-simulator-v2/ .

EXPOSE 8000

CMD ["./neurone-am-simulator-v2"] 
