select 'create database ' || quote_ident(:'database')
where not exists (
  select from pg_database
  where datname = :'database'
)\gexec
