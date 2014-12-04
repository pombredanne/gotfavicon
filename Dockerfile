############################################################
# Dockerfile to run gotfavicon inside a container
# Based on Ubuntu Image
############################################################
FROM ubuntu
MAINTAINER Thomas Sileo
EXPOSE 8000
ADD ./gotfavicon /opt/gotfavicon
ENTRYPOINT /opt/gotfavicon
