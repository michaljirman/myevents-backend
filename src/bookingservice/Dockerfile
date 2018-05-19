# FROM golang:1.9

# WORKDIR /go/src/github.com/michaljirman/myevents-backend
# COPY . .
# WORKDIR src/bookingservice
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bookingservice

FROM debian:jessie
RUN groupadd -r bookingservice && useradd --no-log-init -r -g bookingservice bookingservice

FROM scratch
COPY --from=0 /etc/passwd /etc/passwd
COPY --from=0 /etc/group /etc/group
COPY --chown=bookingservice:bookingservice bookingservice /bookingservice
COPY --chown=bookingservice:bookingservice cert.pem /root/certs/cert.pem
COPY --chown=bookingservice:bookingservice key.pem /etc/ssl/private/key.pem

ENV LISTEN_URL=0.0.0.0:8181
EXPOSE 8181

ENV LISTEN_URL_TLS=0.0.0.0:9292
EXPOSE 9292

CMD [ "/bookingservice" ]