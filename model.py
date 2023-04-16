import collections
import re
import requests
from gensim.models import Doc2Vec
from gensim.models.doc2vec import TaggedDocument
from gensim.utils import simple_preprocess
from sklearn.model_selection import train_test_split


# function to make TaggedDocuments for training. set tags to article name (can also use id)
def create_tagged_document(article):
    return TaggedDocument(words=simple_preprocess(article['content']), tags=article['slug'])

#load json files
SECTIONS = ['news', 'sports', 'opinion']
PAGES = 10

def clean_text(txt):
    clean = re.compile('(<.*?>)|(&nbsp;)|(&amp;)')
    return re.sub(clean, '', txt)

tagged_articles = []


def scrape_data():
    for section in SECTIONS:
       # get max number of pages
       resp = requests.get(f'https://www.thedp.com/section/{section}.json?page=1').json()
       num_pages = resp['pagination']['last']
       
       for page in range(1, num_pages+1):
           resp = requests.get(f'https://www.thedp.com/section/{section}.json?page={page + 1}').json()
           articles = resp['articles']
           print(len(articles))
           for article in articles:
                cleaned_article = clean_text(article['content'])
                print(type(cleaned_article))
                tagged_articles.append(create_tagged_document(cleaned_article))

print(len(tagged_articles))
# split data
train_articles, test_articles = train_test_split(tagged_articles, test_size=0.2)

# train the Doc2Vec - need to tweak vector_size hyperparameter
model = Doc2Vec(vector_size=50, min_count=3, epochs=30)
model.build_vocab(train_articles)
model.train(tagged_articles, total_examples=model.corpus_count, epochs=model.epochs)

# are we just loading the tagged vectors to the db, or should we run like KNN or KMeans ML model?
train_vectors = [model.infer_vector(article.words) for article, _ in train_articles]
test_vectors = [model.infer_vector(article.words) for article, _ in test_articles]

ranks = []
second_ranks = []
for doc_id in range(len(train_vectors)):
    inferred_vector = model.infer_vector(train_vectors[doc_id].words)
    sims = model.dv.most_similar([inferred_vector], topn=len(model.dv))
    rank = [docid for docid, sim in sims].index(doc_id)
    ranks.append(rank)
    second_ranks.append(sims[1])

counter = collections.Counter(ranks)
print(counter)