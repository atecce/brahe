create table artists (

	artist_name varchar(255) not null,

	primary key (artist_name)
);

create table albums (

	artist_name varchar(255) not null,

	album_title varchar(255) not null,
	year  	    year(4)      not null,

	primary key (album_title),
	foreign key (artist_name) references artists (artist_name)
);

create table songs (

	album_title varchar(255) not null,

	track 	    tinyint unsigned not null,
	song_title  varchar(255)     not null,
	
	lyrics text,

	primary key (song_title),
	foreign key (album_title) references albums (album_title)
);
