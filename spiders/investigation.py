import scrapy

class investigation(scrapy.Spider):

	name = 'investigation'

	start_urls = [

		"http://www.lyrics.net/"
	]

	def parse(self, response):

		with open('urls/alphabet_urls.txt', 'w') as f:

			for suburl in response.xpath("//div[@id='page-letter-search']//@href").re("^/artists/[A-Z0]$"): 

				f.write(response.urljoin(suburl + '/99999')+'\n')
