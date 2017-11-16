# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/erezlevip/Ubotrade/WebApplication

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

RUN go get -v github.com/gorilla/mux
RUN go get -v github.com/nu7hatch/gouuid
#RUN go get -v github.com/erezlevip/Ubotrade/...

RUN go get -v github.com/erezlevip/Ubotrade/WebApplication/Handlers"
RUN go get -v github.com/erezlevip/Ubotrade/WebApplication/Middlewares"
RUN go get -v github.com/erezlevip/Ubotrade/Logger"
RUN go get -v github.com/erezlevip/Ubotrade/API/AuthenticationServiceAPI"
RUN go get -v github.com/erezlevip/Ubotrade/API/ServiceAPIFactory"
RUN go get -v github.com/erezlevip/Ubotrade/API/UserServiceAPI"
RUN go get -v github.com/erezlevip/Ubotrade/UserService/Service"
RUN go get -v github.com/erezlevip/Ubotrade/API/BotServiceAPI"
RUN go get -v github.com/erezlevip/Ubotrade/Global"


RUN go install github.com/erezlevip/Ubotrade/WebApplication

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/src/github.com/erezlevip/Ubotrade/WebApplication

# Document that the service listens on port 8080.
EXPOSE 8000