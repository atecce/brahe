# need this for pages
import requests

# need this to parse html
from bs4 import BeautifulSoup

# set url
urls = {
		'http://www.metrolyrics.com',
		'http://www.lyricsfreak.com',
		'http://www.lyrics.net'
       }

# for each url
for url in urls:

	print
	print url
	print

	# get page
	page = requests.get(url)

	# get soup
	soup = BeautifulSoup(page.content, 'lxml')

	for link in soup.find_all('a'):

		print '\t'+str(link.get('href'))
