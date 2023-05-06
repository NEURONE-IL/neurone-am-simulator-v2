#!/bin/bash

mongosh <<EOF
   var cfg = {
        "_id": "rs0",
        "version": 1,
        "members": [
            {
                "_id": 0,
                "host": "localhost:27017",
                "priority": 2
            },
        ]
    };
    rs.initiate(cfg, { force: true });
    //rs.reconfig(cfg, { force: true });
    rs.status();
EOF
sleep 10

mongosh <<EOF
   use admin;
   admin = db.getSiblingDB("admin");
   admin.createUser(
     {
	user: "admin",
        pwd: "admin_password",
        roles: [ { role: "root", db: "admin" } ]
     });
    db.getSiblingDB("admin").auth("admin", "admin_password");
    rs.status();
    use neurone;
    db = db.getSiblingDB("neurone");
    db.createUser(
     {
	user: "neurone",
        pwd: "neur0n3",
        roles: [ { role: "readWrite", db: "neurone" } ]
     });
EOF