#### REASON
I like movies!
I like to watch some movie even twice ot 3 times, and I am not always
finding myself with the possibility to go to the cinema to do it.
So I may or may not download that movie using Torrent, because
I like high-quality video source, and I may or may not be tired of looking
from time to time to check if my desired version is out or not, so...this is a way
to help me and possibly someone else to make this process easier and having fun in doing so!



#### CAPABILITIES
This project aims to create an API capable of returning information about a 
huge collection of movies, along with several available Torrent versions and links.
It will also offer the possibility of subscribing for a specific movie and receiving notifications
about the availability of a specific Torrent version (quality, source, etc.).

This API will be fully developed in [Golang](https://golang.org/).



#### REQUIREMENTS
- Golang 1.11 or grater
- Debian or RedHad based Linux distribution


#### INSTALLATION
Run
```
chmod +x scripts/*
sudo ./scripts/install_environment.sh
sudo ./scripts/install_dependencies.sh
```

#### DOCKER
Docker for this project is available in Goggle Cloud Container Registry under '[eu.gcr.io/gonema/gonema](eu.gcr.io/gonema/gonema)'


#### API ACCESS
Access to this API is available at [https://gonemapi.ruggieri.tech](https://gonemapi.ruggieri.tech)


#### TODO LIST
- [X] Finalize a first version of the API, capable of returning basic information
about the searched movie
- [X] Build Docker image
- [X] Deploy on cloud
- [ ] Build a minimal website version in order to use the API
- [ ] Improve the API. Add Movie information from Imdb, possibly using their API
- [ ] Implement CI (+ Docker integration)
- [ ] Create a local DB using ElasticSearch
- [ ] Improve website with more complex JS
- [ ] Integrate with Slack interactive commands