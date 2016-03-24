from inspiration.song import song

class Poetry:

	songs = [
			song("Black Star", "Mos Def & Talib Kweli Are Black Star", 12, "Thieves in the Night", "inspiration/thieves_in_the_night.txt")
		]

for song in Poetry.songs: print song.lyrics
