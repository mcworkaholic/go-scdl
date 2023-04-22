import music_tag
import json
import os
from datetime import datetime

download_cache = ".\\json\\download-cache.json"
def load(filepath):
    return music_tag.load_file(filepath)
 
# Load the JSON file
with open(download_cache, 'r') as f:
    data = json.load(f)
    for obj in data:
        f = load(obj['file_path'])
        title = f['Title']
        length = f['#length']
        codec = f['#codec']
        channels = f['#channels']
        bitspersample = f['#bitspersample']
        samplerate = f['#samplerate']
        artist = f['artist']
        albumartist = f['albumartist']
        # artwork = f['artwork']
        album = f['album']
        genre = f['genre']
        tracknumber = f['tracknumber']
        year = f['year']
        f['title'] = obj['title']
        f['artist'] = obj['username']
        f['genre'] = obj['genre']
        # Parse the string into a datetime object
        dt = datetime.strptime(obj['created_at'], "%Y-%m-%dT%H:%M:%SZ")
        f['year'] = dt.strftime("%Y-%m-%d")
        f.save()
        # print(title, length, codec, channels, bitspersample, samplerate, albumartist, artist, genre, tracknumber, year)

os.remove(download_cache)
