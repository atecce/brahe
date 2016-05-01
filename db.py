import MySQLdb

from _mysql_exceptions import OperationalError

class canvas:

	# set the canvas
	canvas = MySQLdb.connect('localhost', 'root', charset='utf8')
	brush  = canvas.cursor()

	try: 

		brush.execute("create database canvas")
		canvas.close()

	except: canvas.close()

	def __init__(self):

		# sketch the outline
		canvas, brush = self.prepare()

		brush.execute("""create table if not exists artists (
									  
					name varchar(255) not null, 	  
								  
					primary key (name) 		  
								  
					)""")

		brush.execute("""create table if not exists albums ( 		
										
					title varchar(255) not null, 			
											
					artist_name varchar(255) not null,  		
										
					primary key (title, artist_name), 				
					foreign key (artist_name) references artists (name) 
										
					)""")

		brush.execute("""create table if not exists songs ( 	    	       
									       
					title varchar(255) not null, 	    	       
										       
					album_title varchar(255) not null, 	    	       
									       
					lyrics text, 			    	       
									       
					primary key (album_title, title),
					foreign key (album_title) references albums (title)
					
					)""")

		canvas.close()

	def prepare(self):

		canvas = MySQLdb.connect('localhost', 'root', db='canvas', charset='utf8')
		brush  = canvas.cursor()

		return canvas, brush

	def add_artist(self, artist_name):

		canvas, brush = self.prepare()

		brush.execute("""insert into artists (name)
					values (%s)
					on duplicate key update
					name = name""",
					[artist_name])

		canvas.commit()
		canvas.close()

	def add_album(self, artist_name, album_title):

		canvas, brush = self.prepare()

		brush.execute("""insert into albums (artist_name, title)
					values (%s, %s)
					on duplicate key update
					artist_name = artist_name, title = title""",
					[artist_name, album_title])

		canvas.commit()
		canvas.close()

	def add_song(self, album_title, song_title, lyrics):

		canvas, brush = self.prepare()

		brush.execute("""insert into songs (album_title, title, lyrics)
					values (%s, %s, %s)
					on duplicate key update
					album_title = album_title, title = title, lyrics = lyrics""",
					[album_title, song_title, lyrics])

		canvas.commit()
		canvas.close()

	def get_artists(self):

		canvas, brush = self.prepare()

		brush.execute("select name from artists")

		artists = [item[0] for item in brush.fetchall()]

		canvas.close()

		return artists

	def get_albums(self):

		canvas, brush = self.prepare()

		brush.execute("select title from albums")

		albums = [item[0] for item in brush.fetchall()]

		canvas.close()

		return albums

	def get_songs(self):

		canvas, brush = self.prepare()

		brush.execute("select title from songs")

		songs = [item[0] for item in brush.fetchall()]

		canvas.close()

		return songs

	def get_lyrics(self):

		canvas, brush = self.prepare()

		brush.execute("select lyrics from songs")

		lyrics = [item[0] for item in brush.fetchall()]

		canvas.close()

		return lyrics

#class text:
#
#	import nltk
#
#	canvas = canvas()
#
#	text = MySQLdb.connect('localhost', 'root')
#	mind = text.cursor()
#
#	try: 
#
#		mind.execute("create database text")
#		text.close()
#
#	except: text.close()
#
#	def __init__(self):
#
#		# sketch the outline
#		text, mind = self.prepare()
#
#		mind.execute("""create table if not exists tokens (
#									  
#				       token varchar(255) not null, 	  
#
#				       occurences int not null,
#	  
#				       primary key (token) 	
#				       			  
#				       )""")
#		text.close()
#
#	def prepare(self):
#
#		text = MySQLdb.connect('localhost', 'root', db='text')
#		mind = text.cursor()
#
#		return text, mind
#
#	def set_work(self): 
#
#		# get the work form the canvas
#		canvas, brush = self.canvas.prepare()
#
#		brush.execute("select title, lyrics from songs")
#
#		song_data = ((item[0], item[1]) for item in brush.fetchall())
#
#		canvas.close()
#
#		# interpret it
#		text, mind = self.prepare()
#
#		for song, lyrics in song_data:
#
#			print song
#
#			try:
#
#				mind.execute("""create table if not exists `%s` ( 	    	       
#
#						       id int not null auto_increment,
#										       
#						       token varchar(255) not null,
#
#						       primary key (id)
#						       
#						       )""", [song])
#
#				for token in self.nltk.word_tokenize(lyrics.decode('utf-8')): 
#					
#					mind.execute("""insert into `%s` (token) 
#							       values (%s) 
#							       """, [song, token])
#
#					mind.execute("""insert into tokens (token, occurences)
#							       values (%s, %d)
#							       on duplicate key update 
#							       occurences = occurences+1
#							       """, [token, 1]) 
#
#			except: continue
#
#		text.close()
