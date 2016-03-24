from inspiration.song import song

class Move:

	songs = [
			song("The Foreign Exchange", "Connected", 3, "Raw Life", "inspiration/raw_life.txt")
		]

for song in Move.songs: print song.lyrics
