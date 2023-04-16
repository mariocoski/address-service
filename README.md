# address-service

Simple microservice written in golang to handle CRUD operations on Address domain

# TODO:

- adjust project layout (https://github.com/golang-standards/project-layout)
- setup postgres repository
- dokerize the app
- add authentication middleware 
- add tests
- add http handlers for:
  - POST /api/addresses - to create an address
  - GET /api/addresses - to get all addresses paginated
  - GET /api/addresses/:addressId - to get address by id
  - PATCH /api/addresses/:addressId - to patch address by id
  - DELETE /api/addresses/:addressId - to patch address by id

https://gist.github.com/rhcarvalho/66130d1252d4a7b1fbaeacfe3687eaf3