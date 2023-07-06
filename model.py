import gensim
import os
import json
import logging
from pymongo import MongoClient
from sklearn.model_selection import train_test_split

logging.basicConfig(format='%(asctime)s : %(levelname)s : %(message)s', level=logging.INFO)

def download_data():
    uri = os.environ.get("MONGODB_URI")
    client = MongoClient(uri)
    articles = client.Cluster.articles
    documents = list(articles.find({}, {"slug": 1, "content": 1}))
    training_data, testing_data = train_test_split(documents, test_size=0.2)
    with open("data/training.json", "x") as f:
        json.dump(training_data, f, indent=4)
    with open("data/testing.json", "x") as f:
        json.dump(testing_data, f, indent=4)

def read_corpus(filepath, tokens_only=False):
    with open(filepath) as f:
        articles = json.load(f)
        for article in articles:
            id = int(article["_id"])
            tokens = gensim.utils.simple_preprocess(article["content"])
            yield (tokens if tokens_only else gensim.models.doc2vec.TaggedDocument(tokens, [id]))

def train_model():
    train_corpus = list(read_corpus("data/training.json"))
    assert gensim.models.doc2vec.FAST_VERSION > -1
    model = gensim.models.doc2vec.Doc2Vec(vector_size=50, epochs=20, min_count=10)
    model.build_vocab(train_corpus)
    model.train(train_corpus, total_examples=model.corpus_count, epochs=model.epochs)
    model.save("doc2vec.model")
    return model
