import gensim
import os
import redis
import uvicorn
import subprocess
from datetime import datetime, timedelta
from dbutils import utils
from redis.commands.search.query import Query
from pymongo import MongoClient
from fastapi import FastAPI

app = FastAPI()
model = gensim.models.doc2vec.Doc2Vec.load(
    os.path.join(os.getcwd(), "model", "doc2vec.model")
)
mongo_client = MongoClient(os.environ.get("MONGODB_URI"))
redis_client = redis.Redis().from_url(os.environ.get("REDIS_URI"))
top_k = 5


@app.get("/recommend")
def recommend(slug: str = ""):
    vec = redis_client.hget(f"article:{slug}", "embedding")
    if vec:
        dt = datetime.now() - timedelta(
            days=90
        )  # only show articles published in the last 90 days
        query = (
            Query(
                f"(@timestamp:[({dt.timestamp()} inf])=>[KNN {top_k + 1} @embedding $vec as score]"
            )
            .sort_by("score")
            .return_fields("slug", "headline", "thumbnail_url")
            .dialect(2)
        )
        query_params = {"vec": vec}
        results = redis_client.ft("articles").search(query, query_params).docs
        return results[1:] if results[0]["slug"] == slug else results[:-1]
    return None


if __name__ == "__main__":
    subprocess.Popen(
        [
            "redis-server",
            "--loadmodule",
            "/opt/redis-stack/lib/redisearch.so",
            "--port",
            "6379",
            "--save",
        ]
    )
    utils.refresh_recent_articles()
    utils.initialize_schema()
    uvicorn.run("server:app", port=8080, host="0.0.0.0")
