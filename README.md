# content-recommendation
NLP-Based Content Recommendation Engine for the DP!

## Development

Running the server requires [Python 3](https://www.python.org/downloads/) and [Golang](https://go.dev/dl/).


### Loading the data

If you are using a new MongoDB instance, you can run `go run main.go import` to import articles from the DP API into the instance, and `go run main.go update` to add any recently published articles.

### Setup

You will need Redis Stack running, which can be set up with the following commands:
```
docker pull redis/redis-stack-server
docker run -d --name redis-stack -p 6379:6379 redis/redis-stack-server:latest
```

You will also need to set the following environment variables
- `MONGO_URL`
- `REDIS_URL`

If you are running the server for the first time, run `pipenv install`. Then you can run 
```
pipenv shell
pipenv run python3 server.py --dev --reload
```
to start the server. The `--reload` flag can be omitted if you don't want to (re)load article data from MongoDB into Redis.

## Deployment

First, set the following secrets with:
```
fly secrets set MONGO_URL=...
fly secrets set REDIS_URL=redis://localhost:6379
```

To deploy this service on Fly.io, run `fly deploy --wait-timeout 1200` (the timeout is needed to load the articles into the Redis instance and index them before starting the server).
