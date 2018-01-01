FROM golang:1.9.2-alpine3.6 AS build

# Install tools required to build the project
# We will need to run `docker build --no-cache .` to update those dependencies
RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

# Gopkg.toml and Gopkg.lock lists project dependencies
# These layers will only be re-built when Gopkg files are updated
COPY Gopkg.lock Gopkg.toml /go/src/github.com/stephenneal/go-start/
WORKDIR /go/src/github.com/stephenneal/go-start
# Install library dependencies
RUN dep ensure -vendor-only

# Copy all project and build it
# This layer will be rebuilt when ever a file has changed in the project directory
COPY . .
# Build a binary will all deps to use in scratch
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o /bin/project

# This results in a single layer image
FROM scratch
COPY --from=build /bin/project /bin/project
ENTRYPOINT ["/bin/project"]
#CMD ["--help"]