FROM gliderlabs/alpine:3.4

RUN apk --update add bash curl tar

# Install S6 from static bins
RUN cd / && curl -L https://github.com/just-containers/skaware/releases/download/v1.17.1/s6-eeb0f9098450dbe470fc9b60627d15df62b04239-linux-amd64-bin.tar.gz | tar -xvzf -

# Set up haproxy-api
ADD superside /superside/superside
ADD superside.toml /superside/superside.toml
ADD public /superside/public
ADD docker/s6 /etc

EXPOSE 7779

CMD ["/bin/s6-svscan", "/etc/services"]
