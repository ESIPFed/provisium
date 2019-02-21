# Start from scratch image and add in a precompiled binary
# CGO_ENABLED=0 env go build .
# docker build  --tag="esiplab/provisium:0.0.9"  .
# docker run -d -p 9900:9900  esiplab/provisium:0.0.9
#FROM scratch
FROM alpine

# Add in the datavolume (needs to be mounted from elsewhere)
ADD services /
ADD ./dataVolume/kvStores /dataVolume/kvStores

# Add our binary
CMD ["/sevices"]

# Document that the service listens on this port
EXPOSE 9900
