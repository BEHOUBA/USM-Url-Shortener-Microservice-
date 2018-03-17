# USM-Url-Shortener-Microservice-
Cet api sert a créer des liens urls de redirection correspondant des des vrais url internet: Projet issu de Freecodecamp.org.
# this api allow user to make short version from a long url
run : go build.
and execute the binary (notice that a postgres sql database must be setup with database named: urlsdatabase, and table named: urlmap.
CREATE TABLE URLmap(
ID SERIAL PRIMARY KEY,
ORIGINAL TEXT,
SHORT_URL TEXT);
