from song import song

class Bug:

	songs = [
			song("Broken Social Scene", 	    			    "Forgiveness Rock Record",  11, "Sweetest Kill", 		"sweetest_kill.txt"),
			song("M83", 		    	    			    "Hurry Up, We're Dreaming",  6, "Raconte-Moi Une Histoire", "raconte-moi_une_histoire.txt"),
			song("Radiohead", 	    	    			    "Hail to the Thief",  	 9, "There, There", 		"there_there.txt"),
			song("Daughter", 	    	    			    "If You Leave",  	 	 3, "Youth", 			"youth.txt"),
			song("Handome Boy Modeling School (feat. Roisin & J-Live)", "So.. How's Your Girl",  	 3, "The Truth", 		"the_truth.txt")
		]

for song in Bug.songs:

	print song.lyrics
