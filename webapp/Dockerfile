# Start from scratch image and add in a precompiled binary
# env GOOS=linux go build .
# docker build  --tag="esiplab/provisium:0.0.1"  .
# docker run -d -p 9900:9900  esiplab/provisium:0.0.1
FROM scratch

# Add in the static elements (could also mount these from local filesystem)
ADD webapp /
ADD ./static /static
ADD ./templates /templates
ADD ./dataResources /dataResources

# Add our binary
CMD ["/webapp"]

# Document that the service listens on this port
EXPOSE 9900
