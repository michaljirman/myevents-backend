# FROM golang:1.9

# WORKDIR /go/src/github.com/michaljirman/myevents-backend
# COPY . .
# WORKDIR src/eventservice
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o eventservice

FROM debian:jessie
RUN groupadd -r eventservice && useradd --no-log-init -r -g eventservice eventservice

FROM scratch
COPY --from=0 /etc/passwd /etc/passwd
COPY --from=0 /etc/group /etc/group
COPY --chown=eventservice:eventservice eventservice /eventservice
# COPY --chown=eventservice:eventservice cert.pem /root/certs/cert.pem
# COPY --chown=eventservice:eventservice key.pem /etc/ssl/private/key.pem

ENV LISTEN_URL=0.0.0.0:8181
EXPOSE 8181

# ENV LISTEN_URL_TLS=0.0.0.0:9191
# EXPOSE 9191

USER eventservice

CMD [ "/eventservice" ]