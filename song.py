class song:

	artist = str()
	album  = str()
	track  = int()
	title  = str()
	lyrics = str()

	def __init__(self, artist, album, track, title, filename):

		self.artist = artist
		self.album  = album
		self.track  = track
		self.title  = title

		with open(filename) as f:

			self.lyrics = f.read()
