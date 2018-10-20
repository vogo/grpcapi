
# Create Mongodb in Docker

## 1. Start mongodb

```
cd script
./run-docker-mongodb.sh
```

## 2. Initial mongodb database and oauth2 data

execute `./conn-docker-mongodb.sh` to connect mongodb docker instance, and execute the following commands to create user and initial oauth2 collection data:

```
> use admin
switched to db admin

> db.auth("admin","grpcpass")
1

> use oauth2
switched to db oauth2

> db.createUser(
   {
     user: "oauth2",
     pwd: "oauth2",
     roles: [ "readWrite", "dbAdmin" ]
   }
)
Successfully added user: { "user" : "oauth2", "roles" : [ "readWrite", "dbAdmin" ] }

> db.client.insert({
     "_id" : "000000",
     "secret" : "999999",
     "domain" : "https://localhost",
     "scopes" : [
         "read",
         "manage",
         "admin",
         "view"
     ],
     "grant_types" : [
         "password",
         "refresh_token",
     ]
 })
WriteResult({ "nInserted" : 1 })


> db.user.insert({
    "_id" : "u1",
    "password" : BinData(0, "bTYCH2X8J0a1F9aK1hv//9BYC7/SIyq89OVz3pZvcog="),
    "salt" : BinData(0, "NXl1mO6iV5YpI8YEuy2vTg=="),
    "scopes" : {
        "000000" : "manage,admin"
    }
})
WriteResult({ "nInserted" : 1 })
>
```
