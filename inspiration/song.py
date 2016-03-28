# need these to get and parse html
from bs4 import BeautifulSoup
import requests
import re

# need this for the act of interpretation
import nltk

class song:

	def __init__(self, artist, title, url):

		# set these directly
		self.artist = artist
		self.title  = title

		print
		print '\t\t\t', self.title
		print

		# read the lyrics from a web page
		self.inspire(url)

		# write the lyrics to a text file
		with open('inspiration/' + title, 'w') as f: f.write(self.lyrics)

	def inspire(self, url):

		# set up page
		page = requests.get(url)
		soup = BeautifulSoup(page.content, 'lxml')

		# get lyrics
		if soup.find_all('pre'): 
			
			self.lyrics = soup.find_all('pre')[0].text

			print

			for line in self.lyrics.splitlines():

				print '\t\t\t\t', line

			print

	def interpret(self):

		# tokenize lyrics
		self.tokens = nltk.word_tokenize(self.lyrics)

		# take a distribution of the tokens
		self.token_distribution = nltk.probability.FreqDist(self.tokens)

		# tag parts of speech
		self.pos_tags = nltk.pos_tag(self.tokens)

		# take a distribution of the tags
		self.tag_distribution = nltk.probability.FreqDist(self.pos_tags)
