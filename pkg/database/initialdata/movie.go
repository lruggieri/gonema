package initialdata

type Movie struct{
	Id int `json:"-"`
	Name string `json:"name"`
	Genres []string `json:"genres"`
	ImdbID string `json:"imdb_id"`
	TmdbID string `json:"tmdb_id"`
}
