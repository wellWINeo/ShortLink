# ShortLink
Golang service to cut links with MongoDB and website. Can run standalone, without http-server, but don't support https. 
To run it with nginx i placed placed config in configs dir, also you can find Dockerfile and docker-compose.yml here.

With this service you can make you big and ugly links shorter and prettier. You can share not only links, but also a QR-codes.

# TODO

- [ ] https support
- [ ] connect to MongoDB with login & password
- [ ] fix some issues in layouts
- [ ] history
- [ ] deleting after timeout
