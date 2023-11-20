import argparse
import gensim
import os
import redis
import uvicorn
import subprocess
import numpy as np
from datetime import datetime, timedelta
from dbutils import utils
from redis.commands.search.query import Query
from fastapi import FastAPI
from fastapi_cache import FastAPICache
from fastapi_cache.backends.inmemory import InMemoryBackend
from fastapi_cache.decorator import cache
from pymongo import MongoClient

app = FastAPI()
model = gensim.models.doc2vec.Doc2Vec.load(
    os.path.join(os.getcwd(), "model", "doc2vec.model")
)
redis_client = redis.Redis().from_url(os.environ.get("REDIS_URI"))
mongo_client = MongoClient(os.environ.get("MONGODB_URI"))
top_k = 5


@app.on_event("startup")
async def startup():
    FastAPICache.init(InMemoryBackend())


@app.get("/recommend")
@cache(expire=60 * 60)
def recommend(uuid: str = ""):
    dt = datetime.now() - timedelta(
        days=90
    )  # only show articles published in the last 90 days
    query = (
        Query(
            f"(@timestamp:[({dt.timestamp()} inf])=>[KNN {top_k + 1} @embedding $vec as score]"
        )
        .sort_by("score")
        .return_fields("uuid", "headline", "thumbnail_url")
        .dialect(2)
    )
    key = f"article:{uuid}"
    vec = redis_client.hget(key, "embedding")
    if vec:
        query_params = {"vec": vec}
        results = redis_client.ft("articles").search(query, query_params).docs
        return results[1:] if results[0]["id"] == key else results[:-1]
    doc = mongo_client.Cluster.articles.find_one({"uuid": uuid})
    if doc:
        embedding = [
            float(val)
            for val in list(
                model.infer_vector(gensim.utils.simple_preprocess(doc["content"]))
            )
        ]
        byte_embedding = np.array(embedding, dtype=np.float32).tobytes()
        query_params = {"vec": byte_embedding}
        results = redis_client.ft("articles").search(query, query_params).docs
        return results[1:] if results[0]["id"] == key else results[:-1]
    return None


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Start the backend server")
    parser.add_argument(
        "--dev", action="store_true", help="Run the server in development mode"
    )
    parser.add_argument(
        "--refresh",
        action="store_true",
        help="(Re)load the articles in the Redis instance from MongoDB",
    )
    args = parser.parse_args()
    if not args.dev:
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
    if args.refresh:
        utils.refresh_recent_articles()
        utils.initialize_schema()
    if not args.dev:
        with open("/etc/environment", "w") as f:
            subprocess.run(["env"], stdout=f)
        subprocess.run(["cron"])
    uvicorn.run("server:app", port=8080, host="0.0.0.0")
