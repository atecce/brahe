# need this for pages
import requests

# need this to parse html
from bs4 import BeautifulSoup

# set url
url = 'http://www.metrolyrics.com'

# get page
page = requests.get(url)

# get soup
soup = BeautifulSoup(page.content, 'lxml')

for item in soup.find_all('li'):

	print item
