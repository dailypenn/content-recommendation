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
       print(num_pages)
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
model = Doc2Vec(vector_size=100, min_count=3, epochs=60, min_alpha=0.001, alpha=0.025)
model.build_vocab(tagged_articles)
model.train(tagged_articles, total_examples=model.corpus_count, epochs=model.epochs)

train_vectors = [model.infer_vector(article.words) for article in train_articles]
test_vectors = [model.infer_vector(article.words) for article in test_articles]

# testing
ranks = []
second_ranks = []
# print(len(model.dv))
for idx, doc_id in enumerate(train_articles):
    # print('docid:', doc_id)
    inferred_vector = model.infer_vector(train_articles[idx].words)
    sims = model.dv.most_similar([inferred_vector])
    # print(sims)
    rank = [docid for docid, sim in sims].index(list(train_articles[idx].tags)[0])
    ranks.append(rank)
    second_ranks.append(sims[1])

counter = collections.Counter(ranks)
# print(second_ranks)
counter