import scrapy

class investigation(scrapy.Spider):

	name = 'investigation'

	allowed_domains = 'lyrics.net'

	start_urls = [

		"http://www.lyrics.net/"
	]

	def parse(self, response):

		links = response.xpath("//div[@id='page-letter-search']")

		for link in links: print links
