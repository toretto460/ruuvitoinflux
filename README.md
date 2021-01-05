# Ruuvi to InfluxDB cloud
---

An app to catch all the broadcasted Ruuvi temperature sensors data and send to influxdb

## Usage

To generate the auth token look at the [InfluxDB cloud documentation](https://docs.influxdata.com/influxdb/v2.0/security/tokens/)

```bash
export INFLUX_TOKEN="..."
export INFLUX_BUCKET="..."
export INFLUX_ORG="..."
export INFLUX_HOST="..."

go run cmd/sensor/main.go
```