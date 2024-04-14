sudo docker exec -it mongodb-one mongosh --eval "rs.initiate({_id:'dbrs', members: [{_id:0, host: 'mongodb.one'},{_id:1, host: 'mongodb.two'}]})"
