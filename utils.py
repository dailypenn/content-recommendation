import os
import gensim
import redis
import numpy as np
from redis.commands.json.path import Path
from redis.commands.search.field import TextField, VectorField
from redis.commands.search.indexDefinition import IndexDefinition, IndexType
from pymongo import MongoClient

def initialize_schema():
    redis_client = redis.Redis().from_url(os.environ.get("REDIS_URI"))
    schema = (
        TextField("slug"),
        TextField("headline"),
        TextField("thumbnail_url"),
        VectorField("embedding", algorithm="FLAT", attributes={"TYPE": "FLOAT32", "DIM": 50, "DISTANCE_METRIC": "COSINE"})
    )
    try:
        redis_client.ft("articles").info()
        print("Index 'articles' already exists.")
    except:
        redis_client.ft("articles").create_index(schema, definition=IndexDefinition(prefix=["article:"], index_type=IndexType.HASH))
        print("Succesfully created index 'articles'.")

def refresh_recent_articles():
    model = gensim.models.doc2vec.Doc2Vec.load("doc2vec.model")
    mongo_client = MongoClient(os.environ.get("MONGODB_URI"))
    redis_client = redis.Redis().from_url(os.environ.get("REDIS_URI"))
    pipe = redis_client.pipeline()
    pipe.flushdb()
    cursor = mongo_client.Cluster.articles.find({}, projection={"slug": 1, "headline": 1, "dominantmedia": 1, "content": 1, "createdat": 1}).sort("ctime", -1).limit(1000)
    for doc in cursor:
        embedding = [float(val) for val in list(model.infer_vector(gensim.utils.simple_preprocess(doc['content'])))]
        byte_embedding = np.array(embedding, dtype=np.float32).tobytes()
        key = doc['slug']
        pipe.hset(name=f"article:{key}", mapping={
            "slug": key,
            "headline": doc['headline'],
            "embedding": byte_embedding,
            "thumbnail_url": f'https://snworksceo.imgix.net/dpn/{doc["dominantmedia"]["attachment_uuid"]}.sized-1000x1000.{doc["dominantmedia"]["extension"]}?w=800' 
            if 'dominantmedia' in doc and isinstance(doc['dominantmedia'], dict) 
            else ''
        })
    count = len(pipe) - 1
    print("Error refreshing articles." if False in pipe.execute() else f"Successfully refreshed {count} articles.")
    cursor.close()