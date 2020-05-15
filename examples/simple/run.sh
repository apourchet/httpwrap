#! /bin/bash

set -x

# Bad Creds
curl localhost:3000/pets

# List all 0 pets.
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets

# Get a pet that does not exist.
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets/does_not_exist

# Create pet1 and pet2.
curl -H 'X-PETSTORE-KEY:my-secret-key' -X POST -d '{"name": "pet1", "category": 1}' localhost:3000/pets
curl -H 'X-PETSTORE-KEY:my-secret-key' -X POST -d '{"name": "pet2", "category": 2}' localhost:3000/pets

# List the pets again, should see 2.
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets

# Try to create pet2 again, conflict.
curl -H 'X-PETSTORE-KEY:my-secret-key' -X POST -d '{"name": "pet2", "category": 2}' localhost:3000/pets

# Get pet1 in isolation.
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets/pet1

# Update pet1.
curl -H 'X-PETSTORE-KEY:my-secret-key' -X PUT -d '{"category": 3, "photoUrls": ["photo1"]}' localhost:3000/pets/pet1

# Get the updated pet1.
curl -H 'X-PETSTORE-KEY:my-secret-key' localhost:3000/pets/pet1

# Get a list of pets with photos, should return only pet1.
curl -H 'X-PETSTORE-KEY:my-secret-key' 'localhost:3000/pets/filtered?hasPhotos=true'

# Get a list of pets without photos, should return only pet2.
curl -H 'X-PETSTORE-KEY:my-secret-key' 'localhost:3000/pets/filtered?hasPhotos=false'

# Get a list of pets in categories 1 or 2, should return pet2 only.
curl -H 'X-PETSTORE-KEY:my-secret-key' 'localhost:3000/pets/filtered?categories=1&categories=2'

# Get a list of pets with photos, and in either category 1 or 2, should return empty list.
curl -H 'X-PETSTORE-KEY:my-secret-key' 'localhost:3000/pets/filtered?categories=1&categories=2&hasPhotos=true'

# Clear the petstore for the next run.
curl -H 'X-PETSTORE-KEY:my-secret-key' -X POST localhost:3000/clear
