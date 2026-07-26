# online-subscriptions
An example of simple users online subscriptions aggregation service written on Go

#### Swagger notes
Swagger documentation is integrated into service and is available via endpoint: `/docs/index.html` in your actual service's deployment url, for example: `http://localhost:8080/docs/index.html`  

### Prerequisites of you build host
- Git and make utilities should be installed  
- GoLang compiler and environment should be installed and set up
- Docker with Docker Compose version v2.x.x should be installed and set up

### Build and run
1. Clone project to your working directory
2. In the project's root create `.env` file from project's `.env.example`
3. Customize `.env` file according to your build and deployment environment

#### Build and run service in Docker container
1. Edit project's `docker-compose.yml` file to match your deployment environment
2. Edit `./database/init.sql` that creating service's database name matches `DB_NAME` variable from `.env` file  
3. Run `make docker-rebuild` to compile, build and deploy service to your Docker environment

*Important notes regarding database creation in docker environment*  
*Service's database will be created only on the first run - at the moment when no `online-subscriptions_postgres_data` Docker volume exists and mounted. To re-create the database you should stop Docker containers with `make docker-down` command, then delete the volume by `docker volume rm online-subscriptions_postgres_data`*  

#### Local build and run
- make sure you are able to connect to your PostgreSQL server in your local network
- create PostgreSQL `subscriptions` database
- run `make swagger && make install && make build` to fully build project
- run `make run` to run service locally. SQL migrations will be applied automatically

#### Other useful `make` commands
- `make install` - downloads and installs project dependencies
- `make swagger` - installs swagger and builds swagger docs
- `make build` - compiles the project
- `make test` - runs tests
- `make clean` - clean up project from generated files
- `make docker-build` - builds Docker image
- `make docker-up` - starts service in Docker container
- `make docker-down` - stops service's Docker container
- `make migrate-status` - displays current SQL migrations status
- `make migrate-create` - creates new SQL migration, usage example: `make migrate-create name=create_new_table`
- `make migrate-up` - to apply new SQL migrations
- `make migrate-down` - to revert recent SQL migration
