import scrapy

class investigation(scrapy.Spider):

	name = 'investigation'

	allowed_domains = 'lyrics.net'

	start_urls = [

		"http://www.lyrics.net/"
	]

	def parse(self, response):

		for suburl in response.xpath("//div[@id='page-letter-search']//@href").re("^/artists/[A-Z0]$"): 
			
			print suburl
