import json
import requests
import io
from PIL import Image
import music_tag
from datetime import datetime

def load(filepath):
   return music_tag.load_file(filepath)

    # title = f['Title']
    # length = f['#length']
    # codec = f['#codec']
    # channels = f['#channels']
    # bitspersample = f['#bitspersample']
    # samplerate = f['#samplerate']
    # artist = f['artist']
    # albumartist = f['albumartist']
    # # artwork = f['artwork']
    # album = f['album']
    # genre = f['genre']
    # tracknumber = f['tracknumber']
    # year = f['year']

def create_thumbnail_from_url(url):
    response = requests.get(url)
    img = Image.open(io.BytesIO(response.content))
    img_bytes = io.BytesIO()
    img.save(img_bytes, format='JPEG')
    img_bytes.seek(0)
    return img_bytes.getvalue()

# Load the JSON file
with open(".\\json\\download-cache.json", 'r') as f:
    data = json.load(f)
    for obj in data:
        # load URL from JSON file
        album_artwork_url = obj['artwork_url']
        # Set album artwork from URL
        artwork_data = create_thumbnail_from_url(album_artwork_url)
        f = load(obj['file_path'])
        f['artwork'] = artwork_data
        if (obj['track_format'] == 'single-track'):
            f['album'] = '(Single)'
        f['title'] = obj['title']
        f['artist'] = obj['username']
        f['genre'] = obj['genre']
        # Parse the string into a datetime object
        dt = datetime.strptime(obj['created_at'], "%Y-%m-%dT%H:%M:%SZ")
        year = dt.year
        f['year'] = year
        f.save()
        # print(obj['file_path'] + "\n", obj['file_name']+ "\n", obj['created_at']+ "\n", obj['title']+ "\n", obj['username']+ "\n", obj['genre']+ "\n", obj['artwork_url']+ "\n")

# print(title, length, codec, channels, bitspersample, samplerate, albumartist, artist, genre, tracknumber, year)