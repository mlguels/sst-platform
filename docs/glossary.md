## Terms crossed

## Backend / API

Endpoint: A URL path on your API (e.g., POST /ingest) that handles a request.

Payload: The data you send in a request (often JSON).

Schema: The expected structure of a payload (fields, types, required/optional).

Validation: Checking that incoming data matches the schema.

Telemetry: Is an automated process of collecting, transmitting data to a place to be analyzed

## Database

Row/Record: A single entry in a database table.

Column/Field: A property of that entry (voltage, timestamp).

Primary key: A unique identifier for each row.

Index: A lookup structure to make queries faster (you’ll want one on timestamp).

## Dashboard / UI

Polling: Repeatedly asking the backend for new data every X seconds.

Charting library: Code that draws graphs (Recharts, Chart.js, etc.).

## General

Latency: The time between sending a request and getting a response.

Throughput: How many requests/data points per second your system can handle.

API client: Code in the UI that knows how to call your backend endpoints.
