FROM golang:1.22-alpine

# Set maintainer label: maintainer=[YOUR-EMAIL]
LABEL maintainer="s2310455007@students.fh-hagenberg.at"

# Set working directory: `/src`
WORKDIR /src

# Copy local files to the working directory
COPY app.go ./
COPY go.mod ./
COPY go.sum ./
COPY model.go ./
COPY main.go ./

# List items in the working directory (ls)
RUN ls

# Build the GO app as myapp binary and move it to /usr/
RUN go build -o /usr/cicd-microservices

#Expose port 8888
EXPOSE 8888

# Run the service myapp when a container of this image is launched
CMD ["/usr/cicd-microservices"]