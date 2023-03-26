import json
from gensim.models import Doc2Vec, TaggedDocument
from gensim.utils import simple_preprocess
from sklearn.model_selection import train_test_split

# function to make TaggedDocuments for training. set tags to article name (can also use id)
def create_tagged_document(articles):
    for i, article in enumerate(articles):
        yield TaggedDocument(words=simple_preprocess(article['content']), tags=article['name'])

# Load articles from JSON file (or however we are storing it after webscraping)
filepath = 'articles.json'
with open(filepath, 'r') as f:
    articles = json.load(f)

tagged_articles = list(create_tagged_document(articles))

# split data
train_articles, test_articles = train_test_split(tagged_articles, test_size=0.2)

# train the Doc2Vec - need to tweak vector_size hyperparameter
model = Doc2Vec(vector_size=50, min_count=3, epochs=30)
model.build_vocab(articles)
model.train(tagged_articles, total_examples=model.corpus_count, epochs=model.epochs)

# are we just loading the tagged vectors to the db, or should we run like KNN or KMeans ML model?
train_vectors = [model.infer_vector(article.words) for article, _ in train_articles]
test_vectors = [model.infer_vector(articles.words) for article, _ in test_articles]
