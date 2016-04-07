import scrapy

class meditate(scrapy.Spider):

	name = 'meditate'

	with open("urls/song_urls.txt") as f:

		start_urls = [line.rstrip() for line in f]

	def parse(self, response):

		for lyrics in response.xpath("//pre[@id='lyric-body-text']/text()").extract():

			for line in lyrics.splitlines():

				print line

