# Fetches all open source projects favorited by the user
GET /kudos

# Fetches a favorited open source project by id
GET /kudos/:id

# Creates (or favorites) an open source project for the logged in user
POST /kudos

# Updates a favorited open source project
PUT /kudos/:id

# Deletes (or unfavorites) a favorited open source project
DELETE /kudos/:id