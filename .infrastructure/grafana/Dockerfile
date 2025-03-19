FROM grafana/grafana:latest

ENV PROMETHEUS_HOST='host.docker.internal:9090'

COPY dashboards /etc/grafana/provisioning/dashboards
COPY datasources /etc/grafana/provisioning/datasources

EXPOSE 3000