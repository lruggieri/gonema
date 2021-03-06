#### DESIGN
![Design](design.png)


#### REASON
I like movies!
I like to watch some movie even twice ot 3 times, and I am not always
finding myself with the possibility to go to the cinema to do it.
So I may or may not download that movie using Torrent, because
I like high-quality video source, and I may or may not be tired of looking
from time to time to check if my desired version is out or not, so...this is a way
to help me and possibly someone else to make this process easier and having fun in doing so!



#### CAPABILITIES
This project aims to create a website capable of displaying information about a 
huge collection of movies, along with several available Torrent versions and links.
It will also offer the possibility of subscribing for a specific movie and receiving notifications
about the availability of a specific Torrent version (quality, source, etc.).

This website is developed in [Golang](https://golang.org/) and Javascript.



#### REQUIREMENTS
- Golang 1.12 or grater
- Debian or RedHad based Linux distribution


#### INSTALLATION
Run
```
chmod +x scripts/*
sudo ./scripts/install_environment.sh
sudo ./scripts/install_dependencies.sh
```

#### DOCKER
Docker for this project is available in Google Cloud Container Registry under '[eu.gcr.io/gonema/gonema](eu.gcr.io/gonema/gonemaweb)'
To run the website:
```
docker-compose up -d --build gonemaweb
```

**Note**: if running locally, *docker-compose.override* will take precedence over the standard compose
file. This means that you will have problems with the following environment variables:
* GONEMAES_API_HOST: defines the endpoint to the main resource DB. You either run it locally or
you should delete this variable from the compose.override file, so that the main compose configuration
will take place.


#### TODO LIST
- [X] Finalize a first version of the API, capable of returning basic information
about the searched movie
- [X] Build Docker image
- [X] Deploy on cloud
- [X] Build a minimal website version in order to use the API
- [X] Create a local DB using ElasticSearch
- [X] Improve the API. Add Movie information from Imdb, possibly using their API
- [ ] Implement CI (+ Docker integration)
- [ ] Improve website with more complex JS and CSS
- [ ] Integrate with Slack interactive commands
