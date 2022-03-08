FROM alpine:3

LABEL org.opencontainers.image.source https://github.com/inpublic-io/twitter-api

RUN apk update \
	&& apk -U upgrade \
	&& apk add --no-cache ca-certificates bash \
	&& update-ca-certificates --fresh \
	&& rm -rf /var/cache/apk/*

# adds inpublic user
RUN	addgroup inpublic \
	&& adduser -S inpublic -u 1000 -G inpublic

USER inpublic

COPY --chown=inpublic:inpublic service /usr/local/bin/
RUN chmod +x /usr/local/bin/service

ENTRYPOINT [ "/usr/local/bin/service" ]