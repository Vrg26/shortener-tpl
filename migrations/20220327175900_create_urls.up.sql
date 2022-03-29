CREATE TABLE urls (
                      id bigserial not null primary key,
                      shortUrl varchar(250) unique not null,
                      originUrl varchar (2100) not null,
                      userId bigint,
                      created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX originurl_idx ON urls (originUrl)