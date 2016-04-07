import scrapy

class admiration(scrapy.Spider):

	with open("urls/album_urls.txt", 'w') as f: f.write('')

	name = 'admiration'

	with open("urls/artist_urls.txt") as f:

		start_urls = [line.rstrip() for line in f]

	def parse(self, response):

		with open("urls/album_urls.txt", 'a') as f:

			for suburl in response.xpath("//h3[@class='artist-album-label']//@href").extract(): 

				f.write(response.urljoin(suburl) + '\n')
