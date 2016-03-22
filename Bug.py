from song import song

class Bug:

	# playlist is ordered
	songs = [
			song("Broken Social Scene", 	    			    "Forgiveness Rock Record",  11, "Sweetest Kill", 		"sweetest_kill.txt"),
			song("M83", 		    	    			    "Hurry Up, We're Dreaming",  6, "Raconte-Moi Une Histoire", "raconte-moi_une_histoire.txt"),
			song("Radiohead", 	    	    			    "Hail to the Thief",  	 9, "There, There", 		"there_there.txt"),
			song("Daughter", 	    	    			    "If You Leave",  	 	 3, "Youth", 			"youth.txt"),
			song("Handome Boy Modeling School (feat. Roisin & J-Live)", "So.. How's Your Girl",  	 3, "The Truth", 		"the_truth.txt"),
			song("Eminem (feat. Rihanna)", 				    "Recovery",  	 	15, "Love the Way You Lie", 	"love_the_way_you_lie.txt"),
			song("Elliott Smith", 				    	    "Roman Candle",  	 	15, "No Name #3", 		"no_name_3.txt")
		]

for song in Bug.songs:

	print
	print song.title
	print

	for token in song.token_distribution:

		print '\t'+token, song.token_distribution[token]
