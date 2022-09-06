insert into users (username, password_hash, verification_code) values 
('albatros', '01010101', '22222222');
insert into users (username, password_hash, verification_code) values 
('martin', '01010101', '22222222');

--insert into citedata (user_id, bucket, boltdb) values 
--(1, 'urn:cts:sktlit:skt0001.nyaya002.J1D:', 'urn:cts:sktlit:skt0001.nyaya002.J1D:=>passage1');

--exapmles for use:
--select id, user_id, bucket, boltdb ['urn:cts:sktlit:skt0001.nyaya002.J1D:'] as pass FROM citedata ;
--select id, user_id, bucket, (each(boltdb)).key, (each(boltdb)).value FROM citedata ;
--update citedata set boltdb ['urn:cts:sktlit:skt0001.nyaya002.J1D:']='old passage' where id=1;
--UPDATE citedata SET boltdb = boltdb || hstore('urn:cts:sktlit:skt0001.nyaya002.J1D:1','new passage') where id=1;