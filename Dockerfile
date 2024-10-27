# syntax=docker.io/docker/dockerfile:1.10-labs
# using `labs` in the line above changes "syntax" https://docs.docker.com/build/dockerfile/frontend/
# this syntax allows "exclude" arg for COPY
# dockerfile syntax verisons: https://hub.docker.com/r/docker/dockerfile

# specifies a parent image:
# this image is alpine3.20 + all the stuff you need to build a golang application
# and names this instance 'build'
# cryptic source image names like 'alpine' explained in https://stackoverflow.com/a/59731596/11593686
# alpine versions: https://alpinelinux.org/downloads/
FROM golang:1.23.2-alpine3.20 AS building-image

# mkdir+cd into new directory, we are going to put everything there
WORKDIR /myapp

# copy entire project there (except for what's listed in .dockerignore)
COPY . .

# install all dependencies (of which there are zero, but just as an example, I'll do that anyway)
RUN go mod download

# run go build, name the executable "server.exe" and also disable CGO because people keep telling me that
RUN CGO_ENABLED=0 go build -o server.exe

# switch to a new clean alpine without the golang stuff, much smaller
# General article about so called multi-stage patterns: https://medium.com/swlh/reducing-container-image-size-esp-for-go-applications-db7658e9063a
FROM alpine:3.20 AS running-image

# copy everything from our folder (so, repo + built executable) from our building-image into the same folder but into the second image
# also exclude all the source files, so the final build is even smaller (although it saves like 20kb)
COPY --exclude=**/*.go --exclude=go.mod --from=building-image /myapp /myapp

# notify docker we are going to be using port 8080
EXPOSE 8080

# cd to the folder again
# if this is not done, relative paths inside of the app code will start from root, which is not what we want
WORKDIR /myapp

# tell docker what to run
ENTRYPOINT ["/myapp/server.exe"]
