# rest-to-googlesheet
Simple service that converts a json request into a googlesheet

## Usage

* Checkout this repo 

```
git clone github.com/ggrcha/rest-to-googlesheet

```

* Place client_secret.json on your local cloned dir

* Create/use docker-compose.yml file

```yml
version: "3.7"

services:

    api:
      image: ggrcha/rest-to-googlesheet
      build: .
      ports:
        - 8000:8000
```
* Run "docker-compose up --build"

* You can make a test request using the Postman collection located at postman/ dir