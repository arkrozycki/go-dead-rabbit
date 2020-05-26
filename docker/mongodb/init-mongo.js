db.auth('root', 'password')

db = db.getSiblingDB('go-dead-rabbit')

db.createUser({
  user: 'user',
  pwd: 'password',
  roles: [{
    role: 'dbAdmin',
    db: 'go-dead-rabbit'
  }, {
    role: 'readWrite',
    db: 'go-dead-rabbit'
  }]
})
