package initialdata

import (
	"encoding/csv"
	"errors"
	"github.com/lruggieri/utils"
	"os"
	"path"
	"strconv"
	"strings"
)

type Builder struct{}
func(b *Builder) GetMovies() (oMovies []Movie, oErr error){
	currentFilePath := utils.GetCallerPaths(1)[0]

	moviesFile, err := os.Open(path.Join(currentFilePath,"data","movies.csv"))
	if err != nil{
		return nil,err
	}
	moviesCSVReader := csv.NewReader(moviesFile)
	moviesCSV,err := moviesCSVReader.ReadAll()
	if err != nil{
		return nil,err
	}

	linksFile, err := os.Open(path.Join(currentFilePath,"data","links.csv"))
	if err != nil{
		return nil,err
	}
	linksCSVReader := csv.NewReader(linksFile)
	linksCSV,err := linksCSVReader.ReadAll()
	if err != nil{
		return nil,err
	}


	// Movie table structure: movieId | title | genres
	moviesCSV = moviesCSV[1:] //remove first column (legend)
	maxMoviesLen,err := strconv.Atoi(moviesCSV[len(moviesCSV)-1][0])
	if err != nil{
		return nil,errors.New("cannot find movies max length")
	}
	//initializing to the max possible length
	movies := make([]Movie,maxMoviesLen+1) //important! initialize with len and cap
	for _,movieCSV := range moviesCSV{
		movieID,err := strconv.Atoi(movieCSV[0])
		if err != nil{
			continue
		}
		movieGenres := make([]string,0)
		categoriesStringSplit := strings.Split(movieCSV[2],"|")
		for _, genre := range categoriesStringSplit{
			if genre != "(no genres listed)"{
				movieGenres = append(movieGenres, genre)
			}
		}

		movies[movieID] = Movie{
			Id:     movieID,
			Name:   movieCSV[1],
			Genres: movieGenres,
		}
	}

	//Links table structure: movieId,imdbId,tmdbId
	linksCSV = linksCSV[1:] //remove first column (legend)
	for _,linkCSV := range linksCSV{
		movieID,err := strconv.Atoi(linkCSV[0])
		if err != nil{
			continue
		}
		movies[movieID].ImdbID = linkCSV[1]
		movies[movieID].TmdbID = linkCSV[2]
	}

	return movies,nil
}