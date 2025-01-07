package VideoDatabase

import (
	"database/sql"
	"log"
)

func GetSeriesId(seriesName string) int64 {
	
	var id int64
	
	row := Handle.QueryRow("SELECT * FROM SeriesAPI WHERE Name = $1", seriesName)
	if err := row.Scan(&id, &seriesName); err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return -1
	}
	
	return id
}

func CreateSeries(name string) {
	if _, err := Handle.Exec("INSERT INTO SeriesAPI (Name) VALUES ($1)", name); err != nil {
		log.Fatal(err)
	}
}

func GetVideos(seasonId int64) []Video {
	
	rows, err := Handle.Query(`
	SELECT SeasonInfo.VideoId, Videos.Title, Videos.Size
	FROM SeasonInfo INNER JOIN Videos ON SeasonInfo.VideoId = Videos.Id
	WHERE SeasonId = $1
	ORDER BY EpisodeNumber ASC`, seasonId)
	
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return make([]Video, 0)
	}
	
	videos := make([]Video, 0)
	for rows.Next() {
		var video Video
		if err := rows.Scan(&video.Id, &video.Title, &video.Size); err != nil {
			log.Fatal(err)
		}
		videos = append(videos, video)
	}
	return videos
}

func GetSeasons(seriesId int64) []Season {
	
	rows, err := Handle.Query(`
	SELECT SeriesInfo.SeasonId, SeasonAPI.Name
	FROM SeriesInfo INNER JOIN SeasonAPI ON SeriesInfo.SeasonId = SeasonAPI.Id
	WHERE SeriesId = $1
	ORDER BY SeasonNumber ASC`, seriesId)
	
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return make([]Season, 0)
	}
	
	seasons := make([]Season, 0)
	for rows.Next() {
		var season Season
		if err := rows.Scan(&season.Id, &season.Name); err != nil {
			log.Fatal(err)
		}
		season.Videos = GetVideos(season.Id)
		seasons = append(seasons, season)
	}
	return seasons
}

func GetSeries() []Series {
	
	var allSeries = make([]Series, 0)
	
	// Get all the series
	rows, err := Handle.Query("SELECT * FROM SeriesAPI")
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return make([]Series, 0)
	}
	
	for rows.Next() {
		var series Series
		if err = rows.Scan(&series.Id, &series.Name); err != nil {
			log.Fatal(err)
		}
		series.Tags = GetTags(series.Id)
		series.Seasons = GetSeasons(series.Id)
		allSeries = append(allSeries, series)
	}
	
	return allSeries
}

func GetSeriesNames() []string {
	rows, err := Handle.Query("SELECT * FROM SeriesAPI")
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return make([]string, 0)
	}
	names := make([]string, 0)
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}
	return names
}

func GetSeasonsNames(seriesName string) []string {
	
	var seriesId = GetSeriesId(seriesName)
	
	rows, err := Handle.Query(`
	SELECT SeasonAPI.Name
	FROM SeriesInfo INNER JOIN SeasonAPI ON SeriesInfo.SeasonId = SeasonAPI.Id
	WHERE SeriesId = $1
	ORDER BY SeasonNumber ASC`, seriesId)
	
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return make([]string, 0)
	}
	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		names = append(names, name)
	}
	return names
}