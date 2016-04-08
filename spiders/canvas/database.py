class canvas:

	# set the canvas
	import MySQLdb

	canvas = MySQLdb.connect('localhost', 'root')
	brush  = canvas.cursor()

	try: 

		brush.execute("create database lyrics_net")
		canvas.close()

	except: canvas.close()

	def __init__(self):

		# sketch the outline
		canvas, brush = self.prepare()

		brush.execute("""create table if not exists songs ( 	    	       
										       
					artist varchar(255) not null, 	    	       
									       
					title varchar(255) not null, 	    	       
									       
					lyrics text, 			    	       
									       
					primary key (artist, title)
					
					)""")

		canvas.close()

	def prepare(self):

		canvas = self.MySQLdb.connect('localhost', 'root', db='lyrics_net')
		brush  = canvas.cursor()

		return canvas, brush

	def add_song(self, artist, title, lyrics):

		canvas, brush = self.prepare()

		brush.execute("""insert into songs (artist, title, lyrics)
					values (%s, %s, %s)
					on duplicate key update
					artist = artist, title = title, lyrics = lyrics""",
					[artist, title, lyrics])

		canvas.commit()
		canvas.close()
