# Risk Assessment tool

This is a risk assessment tool!

## Build

* Backend:

Requirement: **[Go](https://go.dev/doc/install)**

```sh
cd backend
make
mkdir ../build
mv webserver ../build/
```

* Frontend:

Requirement: **[npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)** and **[vue/cli](https://cli.vuejs.org/guide/installation.html)**

```sh
cd frontend
npm install
npm run build
mv dist ../build/assets
```

* Run:

1. Must have a MongoDB as the database.  It could be a [MongoDB container](https://hub.docker.com/_/mongo), of course.  I have it with:
   ```sh
   podman run -d -p 27017:27017 --name mongo-db docker.io/library/mongo:latest
   ```
2. Execute:
   ```sh
   cd build
   ./webserver
   ```
   Note: Default MongoDB URI is `mongodb://localhost:27017`.  You can overwrite it with environment variable: `MONGODB_URI`.
3. Launch a browser and go to http://localhost:8080
4. Then, register the first account as an Administrator and use it!

## Some Development Related Things

* Backend: Here is `make test` for unittest
* Frontend: Here is `npm run lint` to run the linter

## Reference

* [Vue CLI](https://cli.vuejs.org/guide/)
* [Gin Web Framework's Documentation](https://gin-gonic.com/docs/)

## Buy Me a Coffee

To Do
