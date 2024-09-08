# tasks_v1
Tasks for a ToDo app

## Running Locally
Spin up local postgres with sample data (requires Docker running):

```
cd local
docker-compose up -d
cd ..
go run ./cmd/tasks -c local/config.json
```
