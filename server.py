import gensim
import os
import redis
import numpy as np
from redis.commands.search.query import Query
from pymongo import MongoClient
from fastapi import FastAPI

app = FastAPI()
model = gensim.models.doc2vec.Doc2Vec.load("doc2vec.model")
mongo_client = MongoClient(os.environ.get("MONGODB_URI"))
redis_client = redis.Redis().from_url(os.environ.get("REDIS_URI"))
top_k = 5
query = (
    Query(f"*=>[KNN {top_k + 1} @embedding $vec as score]")
    .sort_by("score")
    .return_fields("slug", "headline", "thumbnail_url")
    .paging(1, top_k + 1)
    .dialect(2)
)

@app.get("/recommend")
def recommend(slug: str = ''):
    vec = redis_client.hget(f"article:{slug}", "embedding")
    if vec:
        query_params = {"vec": vec}
        results = redis_client.ft("articles").search(query, query_params).docs
        return results
    doc = mongo_client.Cluster.articles.find_one({"slug": slug})
    if doc:
        embedding = [float(val) for val in list(model.infer_vector(gensim.utils.simple_preprocess(doc['content'])))]
        byte_embedding = np.array(embedding, dtype=np.float32).tobytes()
        query_params = {"vec": byte_embedding}
        results = redis_client.ft("articles").search(query, query_params).docs
        return results
    return None
