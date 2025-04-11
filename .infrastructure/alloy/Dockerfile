FROM grafana/alloy:latest

COPY alloy-config.alloy /etc/alloy/alloy-config.alloy

CMD ["run", "--server.http.listen-addr=0.0.0.0:12345", "--storage.path=/var/lib/alloy/data", "/etc/alloy/alloy-config.alloy"]

