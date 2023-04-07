# address-service
Simple microservice written in golang to handle CRUD operations on Address domaingo get -u github.com/go-chi/chi/v5

# TODO:
- adjust project layout (https://github.com/golang-standards/project-layout)
- select db migration tool
- setup postgres repository
- dokerize the app
- add authentication middleware
- add validation library
- add tests
- add http handlers for:
    - POST   /api/addresses - to create an address
    - GET    /api/addresses - to get all addresses paginated
    - GET    /api/addresses/:addressId - to get address by id
    - PATCH  /api/addresses/:addressId - to patch address by id
    - DELETE /api/addresses/:addressId - to patch address by id
