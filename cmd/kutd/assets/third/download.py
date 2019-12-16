from os import path
import os
import urllib.request

urls = [
  'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.15.8/styles/default.min.css',
  'https://cdn.jsdelivr.net/npm/vue@2.6.11/dist/vue.js',
  'https://code.jquery.com/jquery-3.3.1.slim.min.js',
  'https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css',
  'https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css.map',
  'https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js',
  'https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js.map',
  'https://unpkg.com/axios/dist/axios.min.js',
  'https://unpkg.com/axios/dist/axios.min.map',
  'https://unpkg.com/vue-router@3.1.3/dist/vue-router.js',
]

for url in urls:
  fn = url.replace('https://', '')
  d = path.dirname(fn)
  os.makedirs(d, exist_ok=True)
  urllib.request.urlretrieve(url, fn)
