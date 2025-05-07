FROM grafana/grafana:latest

COPY datasources.yml /etc/grafana/provisioning/datasources/
COPY dashboards.yml /etc/grafana/provisioning/dashboards/
COPY dashboards/ /var/lib/grafana/dashboards/

EXPOSE 3000