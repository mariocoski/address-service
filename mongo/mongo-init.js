var db = connect("mongodb://admin:password@localhost:27017/admin");

db = db.getSiblingDB("address-service");

db.createUser({
  user: "user",
  pwd: "password",
  roles: [
    { role: "readWrite", db: "address-service" },
  ],
  passwordDigestor: "server",
});
