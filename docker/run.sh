#!/bin/bash


docker-compose down
docker-compose up -d  
sleep 10
docker-compose exec mongo /usr/bin/mongo -u root -p admin --eval '''if (rs.status()["ok"] == 0) {
    rsconf = {
      _id : "rs0",
      members: [
        { _id : 0, host : "mongo:27017", priority: 1.0 },
      ]
    };
    rs.initiate(rsconf);
}

rs.conf();'''

