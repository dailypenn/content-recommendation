import re
import requests

SECTIONS = ['news', 'sports', 'opinion']
PAGES = 10

def clean_text(txt):
    clean = re.compile('(<.*?>)|(&nbsp;)|(&amp;)')
    return re.sub(clean, '', txt)

def scrape_data():
    for section in SECTIONS:
       for page in range(PAGES):
           resp = requests.get(f'https://www.thedp.com/section/{section}.json?page={page + 1}').json()
           articles = resp['articles']
           for article in articles:
                cleaned_article = clean_text(article['content'])
                slug = article['slug']