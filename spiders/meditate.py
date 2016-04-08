import scrapy

class meditate(scrapy.Spider):

	from canvas.database import canvas

	name = 'meditate'

	with open("urls/unique_song_urls.txt") as f:

		start_urls = [line.rstrip() for line in f]

	canvas = canvas()

	def parse(self, response):

		artist = response.xpath("//h3//a/text()").extract_first()

		title = response.xpath("//h2[@id='lyric-title-text']/text()").extract_first()

		lyrics = response.xpath("//pre[@id='lyric-body-text']/text()").extract_first()

		self.canvas.add_song(artist, title, lyrics)
