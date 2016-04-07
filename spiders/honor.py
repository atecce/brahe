import scrapy

class admiration(scrapy.Spider):

	with open("urls/artist_urls.txt", 'w') as f: f.write('')

	name = 'honor'

	with open("urls/alphabet_urls.txt") as f:

		start_urls = [line.rstrip() for line in f]

	def parse(self, response):

		with open("urls/artist_urls.txt", 'a') as f:

			for suburl in response.xpath("//tr//@href").extract(): 

				f.write(response.urljoin(suburl) + '\n')
