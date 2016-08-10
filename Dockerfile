FROM golang

RUN mkdir /go/src/dns-probe

ADD *.go /go/src/dns-probe/
RUN go get dns-probe
RUN go install dns-probe

ENV PORT 7777
ENV INTERFACE eth0

EXPOSE ${PORT}
CMD /go/bin/dns-probe
