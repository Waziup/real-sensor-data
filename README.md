# Sensor-data Simulator

This WebApp extracts sensor data from thingspeak public API, organizes it and can be configured to push data to the Waziup platform for testing, simulation and AI development

To see details and examples about the APIs head over to the [API Documentation](API.md).

## Quick Setup

To setup and run the WebApp simply copy and execute the following code into your terminal:

```
git clone https://github.com/Waziup/sensor-data-simulator.git
cd sensor-data-simulator
sudo docker-compose up -d
```

Note: _The database, tables, indices, etc will be created automatically on the first launch and will be filled with some random data._

Once it is up, you can see the status of moving scooters in your browser: http://localhost:8080/

## Development

To activate the development mode you need to open `docker-compose.yml` file, under the desired service (e.g. server), change the target to development:

```
    build:
      context: ./
      target: development  # development | test | production (default)
```

## ENV variables

- `SERVING_ADDR`: Service address for the API server and the UI
- `DATA_EXTRACTION_INTERVAL`: The data extraction interval (Default is 60 minutes).
- `WAZIUP_API_PATH`: Waziup API Path

- `POSTGRES_DB`: PostgreSQL database name
- `POSTGRES_USER`: PostgreSQL username with correct authorizations
- `POSTGRES_PASSWORD`: PostgreSQL password
- `POSTGRES_PORT`: PostgreSQL port
- `POSTGRES_HOST`: PostgreSQL Hostname
