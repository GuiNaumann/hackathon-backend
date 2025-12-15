FROM scratch

COPY . /app/

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/app/server"]
