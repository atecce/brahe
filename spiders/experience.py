import scrapy

class experience(scrapy.Spider):

	with open("urls/song_urls.txt", 'w') as f: f.write('')

	name = 'experience'

	with open("urls/album_urls.txt") as f:

		start_urls = [line.rstrip() for line in f]

	def parse(self, response):

		with open("urls/song_urls.txt", 'a') as f:

			for suburl in response.xpath("//strong//@href").re("^/lyric/.*$"): 

				f.write(response.urljoin(suburl) + '\n')
