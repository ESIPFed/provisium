# Start from scratch image and add in a precompiled binary
# env GOOS=linux go build .
# docker build  --tag="esiplab/provohash:0.0.1"  .
# docker run -d -p 9900:9900  esiplab/provohash:0.0.1
#FROM scratch
FROM alpine

# Add in the static elements (could also mount these from local filesystem)
ADD provohash /
ADD ./kvStores /kvStores
ADD ./static /static
ADD ./templates /templates

# Add our binary
CMD ["/provohash"]

# Document that the service listens on this port
EXPOSE 9911
