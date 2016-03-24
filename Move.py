from inspiration.song import song

class Move:

	songs = [
			song("The Foreign Exchange", "Connected", 	       3, "Raw Life", 	       "inspiration/raw_life.txt"),
			song("Broken Social Scene",  "To Be You and Me",       6, "Major Label Debut", "inspiration/major_label_debut.txt"),
			song("Lupe Fiasco",  	     "Lupe Fiasco's The Cool", 5, "Superstar", 	       "inspiration/superstar.txt")
		]

for song in Move.songs: print song.lyrics
