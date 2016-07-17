#!/usr/bin/env perl6

use DBIish;

my $dbh = DBIish.connect("SQLite", :database<lyrics_net.db>);

my $sth = $dbh.prepare(q:to/STATEMENT/);
	SELECT * 
	FROM songs
	STATEMENT

$sth.execute();

for $sth.allrows() { .print }

$sth.finish;

$dbh.dispose;
