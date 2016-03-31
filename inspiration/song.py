# need these to get and parse html
from bs4 import BeautifulSoup
import requests
import re

# need this for the act of interpretation
import nltk

# need this to create the canvas
import os

class song:

	artist = str()
	title  = str()
	lyrics = str()

	def __init__(self, artist, title, url):

		# set these directly
		self.artist = artist
		self.title  = title

		# set filename
		filename = 'inspiration/' + artist + '/' + title

		# check is we already recorded the song
		if not os.path.isfile(filename):

			print
			print '\t\t\t', self.title
			print

			# read the lyrics from a web page
			self.inspire(url)

			# check for blank canvas
			if not os.path.exists(os.path.dirname(filename)): os.makedirs(os.path.dirname(filename))

			# write the lyrics to a text file
			try:

				# write the lyrics to the file
				with open(filename, 'w') as f: f.write(self.lyrics)

			# hacky, introduced because of Pink
			except IOError: pass

	def inspire(self, url):

		# set up page
		page = requests.get(url)
		soup = BeautifulSoup(page.content, 'lxml')

		# get lyrics
		if soup.find_all('pre'): 
			
			self.lyrics = soup.find_all('pre')[0].text.encode('ascii', 'ignore')

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
