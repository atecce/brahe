create table artists (

	name varchar(255) not null,

	primary key (name)
);

create table albums (

	artist_name varchar(255) not null,

	title varchar(255) not null,
	year  year(4)      not null,

	primary key (title),
	foreign key (artist_name) references artists (name)
);

create table songs (

	id int unsigned not null auto_increment,

	album_title varchar(255) not null,

	track  tinyint unsigned not null,
	title  varchar(255)     not null,
	length time,
	
	lyrics text,

	primary key (id),
	foreign key (album_title) references albums (title)
);
