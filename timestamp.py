# this hack is embarassing

def time_width(unit):

	if   len(unit) == 0: unit = '00'
	elif len(unit) == 1: unit = '0' + unit

	return unit

def song_length_to_time(song_length):

	time_segments = song_length.split(':')

	length = len(time_segments)

	count = int()

	time = list()

	hours   = str()
	minutes = str()
	seconds = str()

	for index in range(length):

		if   count == 0: seconds = time_segments[length-index-1]
		elif count == 1: minutes = time_segments[length-index-1]
		elif count == 2: hours   = time_segments[length-index-1]

		count += 1

	if   len(seconds) == 0: seconds = '00'
	elif len(seconds) == 1: seconds = '0' + seconds

	if   len(minutes) == 0: minutes = '00'
	elif len(minutes) == 1: minutes = '0' + minutes

	if   len(hours) == 0: hours = '00'
	elif len(hours) == 1: hours = '0' + hours

	time = hours + ':' + minutes + ':' + seconds

	return time
