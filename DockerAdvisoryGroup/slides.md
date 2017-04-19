layout: true

---
class: middle, center

# DAG ⚡️ talk - _Build_

_Dave Henderson (@hairyhenderson)_
_2016-04-22_

---

## Background

- I work at Qlik, on qlikcloud.com
- Helped to take the site from 0 to fully-Dockerized Linux services in my first 8 months
  - _(some non-Dockerized Windows services remain)_
- Now breaking up the monoliths into microservices
  - ~30 Dockerized services behind qlikcloud.com, plus a handful of "associated" images (tools, testing things, etc)
- most services on a home-grown poor-man's orchestration system
  - but Swarm is now deployed (using Docker for AWS) and the first services are starting to get deployed soon
- Other parts of the company (largely windows-based Enterprise products, and newer mobile products) starting to adopt Docker too
- 99% of the devs and ops behind qlikcloud.com run macOS and deploy Linux services
  - but almost everyone else in R&D runs Windows (mostly old Win10 versions, some Win7 still)

---

## Build

_Borja's suggested topics:_

- how do you build your applications and artifacts?
- do you rely solely on docker build or do you have some custom script/makefile?
- are these scripts portable?
- do you use some other tooling?

---

### The ways we build Docker images

- single `Dockerfile`, just a `docker build -t whatever .`
- single `Dockerfile`, but with a `Makefile`
- dual `Dockerfile`, with `Makefile`

---

#### vanilla `docker build`

- fairly rare
- most of these end up as Automated Builds on DockerHub
- tend to be the _very_ simple services, or utilities

---

#### `Makefile` with one `Dockerfile`

- majority use this approach
- `build` target is mostly just:
  ```makefile
  build:
      docker build -t $(DOCKER_REPO):$(VER) .
  ```
- why a `Makefile`?
  - `make build`/`make test`/etc... are all very familiar to most of our devs
    - and easier to remember than
    ```console
    $ docker build -t mynamespace/somelong-imagename:whatwasthatversionagain .
    ```
  - computing the version (tag) is sometimes more than just a pass-through from the environment
      ```makefile
      build:
          docker build -t $(DOCKER_REPO):$(shell cat $(VERSION_LOC)/version.txt)
      ```

---

#### The two-step build

- used by a number of our services, which depend on private modules -- need credentials at build time
- two `Dockerfile`s:
  - `Dockerfile.build`
```dockerfile
FROM something
COPY super-sensitive-file-with-secrets ./
RUN the_build
RUN rm ./super-sensitive-file-with-secrets
RUN tar cfz app.tar.gz * .??*
```
  - `Dockerfile`
```dockerfile
FROM something
ADD app.tar.gz ./
```
- The `Makefile` target generally looks like:
```Makefile
build:
    cp ~/super-sensitive-file-with-secrets .
    docker build -t build-image -f Dockerfile.build .
    docker create --name build-container build-image
    docker cp build-container:/usr/src/app/app.tar.gz .
    docker build -t $(DOCKER_REPO):$(VER) .
```
- `docker build --squash` or `COPY --from` might help simplify these...
