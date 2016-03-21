#
# I should not like my writing to spare other people the trouble of thinking. 
# But, if possible, to stimulate someone to thoughts of their own.
#

# need this for pages
import requests

# need this to parse html
from bs4 import BeautifulSoup

# set url
url = 'http://www.lyrics.net'

# get page
page = requests.get(url)

# get soup
soup = BeautifulSoup(page.content, 'lxml')

# for each link
for link in soup.find_all('a'):

	# get the suburl
	suburl = link.get('href')

	# check for artist subpages
	if 'artists' in suburl:

		# hacky
		if '.php' in suburl or 'http' in suburl: continue

		# set artist url
		artist_url = url + suburl

		# get artist page
		artist_page = requests.get(artist_url)

		# get artist soup
		artist_soup = BeautifulSoup(artist_page.content, 'lxml')

		print
		print suburl
		print

		# get artist links
		for artist_link in artist_soup.find_all('a'):

			print '\t', artist_link.get('href')
