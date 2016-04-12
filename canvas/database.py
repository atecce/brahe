class canvas:

	# set the canvas
	import MySQLdb

	import nltk

	canvas = MySQLdb.connect('localhost', 'root')
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

		canvas = self.MySQLdb.connect('localhost', 'root', db='canvas')
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

	def set_work(self): 

		canvas, brush = self.prepare()

		brush.execute("select title, lyrics from songs")

		song_data = ((item[0], item[1]) for item in brush.fetchall())

		for song, lyrics in song_data:

			print song

			brush.execute("""create table if not exists `%s` ( 	    	       
										       
						token varchar(255) not null,

						primary key (token)
						
						)""", [song])

			for token in self.nltk.word_tokenize(lyrics.decode('utf-8')): 
				
				brush.execute("""insert into `%s` (token) 
							values (%s) 
							on duplicate key update
							token = token
							""", [song, token])

		canvas.close()
