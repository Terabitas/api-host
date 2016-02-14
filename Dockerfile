FROM centurylink/ca-certs

COPY ./api-host /api-host
COPY ./io.api-host.conf /api-host.conf

CMD ["./api-host", "--config=./api-host.conf"]

EXPOSE 8080