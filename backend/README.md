Backend Folder structure

cmd/server - entrypoint of application
internal
/ingestion - recieves and normalize telemetry - Takes in data then sends it to a place for it to be reviewed
/anylitics - threshold rules & checks - Gives the data boundaries which will trigger an alert
/control - command dispatch - Sends commands
/storage - db access layer - Connection to the database
/tests - test/simulation - Tests how components of the app work together

pkg
/api - route/handlers DTOs - How the front-end / devices talk to the backend and data transfer objects representing / enforcing the data passed in or out
/models - shared structs: Telemetry, Command, Alert - Where shared components used across the backend live
/utils - logging, configs, errors
