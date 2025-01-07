package VideoDatabase

import (
	"database/sql"
	"log"
)

func GetTagId(tagName string) int64 {
	
	var id int64
	
	row := Handle.QueryRow("SELECT * FROM TagsAPI WHERE Name = $1", tagName)
	if err := row.Scan(&id, &tagName); err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return -1
	}
	
	return id
}

func CreateTag(tagName string) Tag {
	
	tag := Tag{ Id: GetTagId(tagName), Name: tagName }
	
	if tag.Id != -1 {
		return tag
	}
	
	var result sql.Result
	var err error
	
	if result, err = Handle.Exec("INSERT INTO TagsAPI (Name) VALUES ($1)", tag); err != nil {
		log.Fatal(err)
	}
	
	if tag.Id, err = result.LastInsertId(); err != nil {
		log.Fatal(err)
	}
	
	return tag
}

func AddTagByName(seriesid int64, tagName string) {
	
	tag := CreateTag(tagName)
	
	row := Handle.QueryRow("SELECT * FROM TagsInfo WHERE SeriesId = $1 AND TagId = $2", seriesid, tag.Id)
	if row.Err() != sql.ErrNoRows {
		return
	}
	if _, err := Handle.Exec("INSERT INTO TagsInfo VALUES ($1, $2)", seriesid, tag.Id); err != nil {
		log.Fatal(err)
	}
}

func GetTags(seriesid int64) []Tag {
	
	rows, err := Handle.Query(`
		SELECT TagsAPI.Id, TagsAPI.Name
		FROM TagsAPI
		INNER JOIN TagsInfo ON TagsAPI.Id = TagsInfo.TagId
		WHERE TagsInfo.SeriesId = $1`,
	seriesid)
	
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return make([]Tag, 0)
	}
	
	var allTags = make([]Tag, 0, 1)
	
	for rows.Next() {
		tag := Tag{}
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			log.Fatal(err)
		}
		allTags = append(allTags, tag)
	}
	
	return allTags
}