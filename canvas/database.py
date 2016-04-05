class canvas:

	# set the canvas
	import MySQLdb

	canvas = MySQLdb.connect('localhost', 'root')
	brush  = canvas.cursor()

	try: 

		brush.execute("create database lyrics_net")
		canvas.close()

	except: canvas.close()

	# sketch the outline
	canvas = MySQLdb.connect('localhost', 'root', db ='lyrics_net')
	brush  = canvas.cursor()

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

	def draw(self, query, args):

		# sketch the outline
		canvas = self.MySQLdb.connect('localhost', 'root', db='lyrics_net')
		brush  = canvas.cursor()

		brush.execute(query, args)
		canvas.commit()

		canvas.close()

	def caught_up(self, song, album):

		# compare yourself to others
		canvas = self.MySQLdb.connect('localhost', 'root', db='lyrics_net')
		brush  = canvas.cursor()

		if brush.execute("select title, album_title from songs where title=%s and album_title=%s", (song, album)):
			
			canvas.close()
			return True

		canvas.close()
