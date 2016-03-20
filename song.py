import nltk

class song:

	# each song has the following attributes
	artist = str()
	album  = str()
	track  = int()
	title  = str()
	lyrics = str()

	def __init__(self, artist, album, track, title, filename):

		# set these directly
		self.artist = artist
		self.album  = album
		self.track  = track
		self.title  = title

		# read the lyrics from a text file
		with open(filename) as f:

			self.lyrics = f.read()

		# tokenize lyrics
		self.tokens = nltk.word_tokenize(self.lyrics)

		# take a distribution of the tokens
		self.token_distribution = nltk.probability.FreqDist(self.tokens)

		# tag parts of speech
		self.pos_tags = nltk.pos_tag(self.tokens)

		# take a distribution of the tags
		self.tag_distribution = nltk.probability.FreqDist(self.pos_tags)
