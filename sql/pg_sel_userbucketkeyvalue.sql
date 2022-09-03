SELECT value FROM
  (SELECT  skeys(dict) key, svals(dict) value FROM citedata where user_id=($1) and bucket=($2)) AS stat
  where key= ($3);