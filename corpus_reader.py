import sqlite3
from nltk.tokenize import word_tokenize, sent_tokenize

def prepare():

	# connect to database
	con = sqlite3.connect("lyrics_net")
	cur = con.cursor()

	# return the connection
	return con, cur

def artists():

	# get the connection
	con, cur = prepare()

	# execute the query
	cur.execute("select name from artists")
	artists = (item[0] for item in cur.fetchall())

	# close the connection
	con.close()

	return artists

def albums(artist=None):

	# get the connection
	con, cur = prepare()

	# set the query
	query = "select title from albums"

	if artist: query += " where artist_name=\"%s\";" % artist

	# execute the query
	cur.execute(query)
	albums = (item[0] for item in cur.fetchall())

	# close the connection
	con.close()

	return albums

def songs(artist=None, album=None):

	# get the connection
	con, cur = prepare()

	# set the query TODO exploits here
	query = "select title from songs"

	if   album:  query += " where album_title=\"%s\";" % album
	elif artist: query += " where album_title in (select title from albums where artist_name=\"%s\");" % artist

	# execute the query
	cur.execute(query)
	songs = (item[0] for item in cur.fetchall())

	# close the connection
	con.close()

	return songs

def lyrics(artist=None, album=None, song=None):

	# get the connection
	con, cur = prepare()

	# set the query TODO exploits here
	query = "select lyrics from songs"

	if   song:   query += " where title=\"%s\";" % song
	elif album:  query += " where album_title=\"%s\";" % album
	elif artist: query += " where album_title in (select title from albums where artist_name=\"%s\");" % artist

	# execute the query
	cur.execute(query)
	text = (item[0] for item in cur.fetchall())

	# close the connection
	con.close()

	return text

def words(artist=None, album=None, song=None):

	# initialize list of words
	word_list = list()

	count = int()

	# tokenize all the text
	for text in lyrics(artist, album, song): word_list += word_tokenize(text)

	return word_list

def sents(artist=None, album=None, song=None):

	# initialize list of sentences
	sent_list = list()

	# tokenize all the text
	for text in lyrics(artist, album, song): sent_list += sent_tokenize(text)

	return sent_list

def paras(): pass
def tagged_words(): pass
def tagged_sents(): pass
def tagged_paras(): pass
def chunked_sents(): pass
def parsed_sents(): pass
def parsed_paras(): pass
def xml(): pass
def raw(): pass
