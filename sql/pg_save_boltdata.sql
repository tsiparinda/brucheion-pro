UPDATE citedata SET boltdb = case when boltdb is null then hstore(($3),($4))
else boltdb || hstore(($3), ($4)) end where user_id=($1) and bucket=($2);