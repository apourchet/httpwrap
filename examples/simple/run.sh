#! /bin/bash

set -x

curl localhost:3000/pets
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets/does_not_exist
curl -H 'X-PETSTORE-KEY:my-secret-key' -X POST -d '{"name": "pet1", "category": 1}' localhost:3000/pets
curl -H 'X-PETSTORE-KEY:my-secret-key' -X POST -d '{"name": "pet2", "category": 2}' localhost:3000/pets
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets
curl -H 'X-PETSTORE-KEY:my-secret-key' -X POST -d '{"name": "pet2", "category": 2}' localhost:3000/pets
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets/pet1
curl -H 'X-PETSTORE-KEY:my-secret-key' -X PUT -d '{"category": 3, "photoUrls": ["photo1"]}' localhost:3000/pets/pet1
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets/pet1
curl -H 'X-PETSTORE-KEY:my-secret-key' 'localhost:3000/pets/filtered?hasPhotos=true'
curl -H 'X-PETSTORE-KEY:my-secret-key' 'localhost:3000/pets/filtered?hasPhotos=false'
curl -H 'X-PETSTORE-KEY:my-secret-key' 'localhost:3000/pets/filtered?categories=1&categories=2'
curl -H 'X-PETSTORE-KEY:my-secret-key' 'localhost:3000/pets/filtered?categories=1&categories=2&hasPhotos=true'

# Clear the petstore for the next run.
curl -H 'X-PETSTORE-KEY:my-secret-key' -X POST localhost:3000/clear
