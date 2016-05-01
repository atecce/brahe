from db import canvas

from nltk.tokenize import word_tokenize, sent_tokenize

canvas = canvas()

class reader:

	def words(self):

		word_list = list()

		count = int()

		for lyrics in canvas.get_lyrics(): 

			count += 1

			print count

			word_list += word_tokenize(lyrics)

		return word_list

	def sents(self):

		sent_list = list()

		count = int()

		for lyrics in canvas.get_lyrics(): 

			count += 1

			print count

			sent_list += sent_tokenize(lyrics)

		return sent_list

	def paras(self): pass
	def tagged_words(self): pass
	def tagged_sents(self): pass
	def tagged_paras(self): pass
	def chunked_sents(self): pass
	def parsed_sents(self): pass
	def parsed_paras(self): pass
	def xml(self): pass
	def raw(self): pass

test = reader()

print len(test.sents())
