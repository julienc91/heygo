package database

import (
    "errors"
)

type Video struct {
    Id int
    Title string
    Path string
    Slug string
}

func GetVideoFromSlug(slug string) (Video, error) {

    stmt, err := db.Prepare(`SELECT id, title, path, slug
FROM videos
WHERE slug = ?;`)
    if err != nil {
        return Video{}, err
    }

    var video Video
    err = stmt.QueryRow(slug).Scan(&video.Id, &video.Title, &video.Path,
        &video.Slug)
    if err != nil {
        return Video{}, errors.New("Unknown slug")
    }

    return video, nil
}

func GetAllVideos() ([]Video, error) {

    stmt, err := db.Prepare(`SELECT id, title, path, slug FROM videos;`)
    if err != nil {
        return nil, err
    }

    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []Video

    for rows.Next() {
        var video Video
        if err := rows.Scan(&video.Id, &video.Title, &video.Path,
            &video.Slug); err != nil {
            return nil, err
        }

        result = append(result, video)
    }
    
    return result, nil
}

func UpdateVideo(videoId int, values map[string]interface{}) error {

    return UpdateRow(videoId, values, []string{"path", "slug", "title"}, "videos")
}

func InsertVideo(values map[string]interface{}) (int, error) {

    return InsertRow(values, []string{"title", "path", "slug"}, "videos")
}
