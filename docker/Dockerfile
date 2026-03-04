FROM cosmtrek/air:v1.64.5

RUN go install github.com/go-delve/delve/cmd/dlv@latest

ENTRYPOINT ["docker/entrypoint.sh"]
