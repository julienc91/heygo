package database

import (
	"errors"
	"github.com/julienc91/heygo/globals"
	"strings"
)

func getVideoGroupsFromVideoId(id int64) ([]globals.VideoGroup, error) {

	var query = "SELECT video_groups.id, video_groups.title FROM video_groups "
	query += "INNER JOIN video_classification ON video_classification.videos_id=? AND video_groups.id = video_classification.video_groups_id;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.VideoGroup{} })
	if err != nil {
		return nil, err
	}
	var videoGroups []globals.VideoGroup
	for _, v := range res {
		videoGroups = append(videoGroups, v.(globals.VideoGroup))
	}
	return videoGroups, err
}

func getVideoGroupFromId(id int64) (globals.VideoGroup, error) {

	var query = "SELECT id, title FROM video_groups WHERE id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.VideoGroup{} })
	if err != nil {
		return globals.VideoGroup{}, err
	} else if len(res) != 1 {
		return globals.VideoGroup{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.VideoGroup), nil
}

func getVideosFromVideoGroupId(id int64) ([]globals.Video, error) {

	var query = "SELECT videos.id, videos.title, videos.slug, videos.path, videos.imdb_id FROM videos "
	query += "INNER JOIN video_classification ON video_classification.videos_id = videos.id AND video_classification.video_groups_id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.Video{} })
	if err != nil {
		return nil, err
	}
	var videos []globals.Video
	for _, v := range res {
		videos = append(videos, v.(globals.Video))
	}
	return videos, nil
}

func getVideoFromId(id int64) (globals.Video, error) {

	var query = "SELECT id, title, slug, path, imdb_id FROM videos WHERE id=?;"
	var params = []interface{}{id}
	res, err := getDb(query, params, func() interface{} { return globals.Video{} })
	if err != nil {
		return globals.Video{}, err
	} else if len(res) != 1 {
		return globals.Video{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.Video), nil
}

func getVideoFromSlug(slug string) (globals.Video, error) {

	var query = "SELECT id, title, slug, path, imdb_id FROM videos WHERE slug=?;"
	var params = []interface{}{slug}
	res, err := getDb(query, params, func() interface{} { return globals.Video{} })
	if err != nil {
		return globals.Video{}, err
	} else if len(res) != 1 {
		return globals.Video{}, errors.New("Unexpected result from database")
	}
	return res[0].(globals.Video), nil
}

func getAllVideos() ([]globals.Video, error) {

	var query = "SELECT id, title, slug, path, imdb_id FROM videos;"
	res, err := getDb(query, nil, func() interface{} { return globals.Video{} })
	if err != nil {
		return nil, err
	}

	var videos []globals.Video
	for _, v := range res {
		videos = append(videos, v.(globals.Video))
	}
	return videos, err
}

func getAllVideoGroups() ([]globals.VideoGroup, error) {

	var query = "SELECT id, title FROM video_groups;"
	res, err := getDb(query, nil, func() interface{} { return globals.VideoGroup{} })
	if err != nil {
		return nil, err
	}

	var videoGroups []globals.VideoGroup
	for _, v := range res {
		videoGroups = append(videoGroups, v.(globals.VideoGroup))
	}
	return videoGroups, err
}

func updateVideoClassificationFromVideoId(id int64, videoGroups []globals.VideoGroup) ([]globals.VideoGroup, error) {

	if err := deleteVideoClassificationFromVideoId(id); err != nil {
		return nil, err
	}

	var query = "INSERT INTO video_classification (videos_id, video_groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, videoGroup := range videoGroups {
		values = append(values, "(?, ?)")
		params = append(params, id, videoGroup.Id)
	}
	query += strings.Join(values, ", ") + ";"

	if err := insertDb(query, params); err != nil {
		return nil, err
	}
	return getVideoGroupsFromVideoId(id)
}

func updateVideoClassificationFromVideoGroupId(id int64, videos []globals.Video) ([]globals.Video, error) {

	if err := deleteVideoClassificationFromVideoGroupId(id); err != nil {
		return nil, err
	}

	var query = "INSERT INTO video_classification (videos_id, video_groups_id) VALUES "
	var params = []interface{}{}
	var values []string
	for _, video := range videos {
		values = append(values, "(?, ?)")
		params = append(params, video.Id, id)
	}
	query += strings.Join(values, ", ") + ";"

	if err := insertDb(query, params); err != nil {
		return nil, err
	}

	return getVideosFromVideoGroupId(id)
}

func updateVideoGroup(videoGroup globals.VideoGroup) (globals.VideoGroup, error) {

	var query = "UPDATE video_groups SET title=? WHERE id=?;"
	var params = []interface{}{videoGroup.Title, videoGroup.Id}
	if err := updateDb(query, params); err != nil {
		return globals.VideoGroup{}, err
	}

	return getVideoGroupFromId(videoGroup.Id)
}

func updateVideo(video globals.Video) (globals.Video, error) {

	var query = "UPDATE videos SET title=?, path=?, imdb_id=?, slug=? WHERE id=?;"
	var params = []interface{}{video.Title, video.Path, video.ImdbId, video.Slug, video.Id}
	if err := updateDb(query, params); err != nil {
		return globals.Video{}, err
	}

	return getVideoFromId(video.Id)
}

func insertVideoGroup(videoGroup globals.VideoGroup) (globals.VideoGroup, error) {

	var query = "INSERT INTO video_groups (title) VALUES (?);"
	var params = []interface{}{videoGroup.Title}
	id, err := insertAndGetId(query, params)
	if err != nil {
		return globals.VideoGroup{}, err
	}

	return getVideoGroupFromId(id)
}

func insertVideo(video globals.Video) (globals.Video, error) {

	var query = "INSERT INTO videos (title, path, slug, imdb_id) VALUES (?, ?, ?, ?);"
	var params = []interface{}{video.Title, video.Path, video.Slug, video.ImdbId}
	id, err := insertAndGetId(query, params)
	if err != nil {
		return globals.Video{}, err
	}

	return getVideoFromId(id)
}

func deleteVideoClassificationFromVideoId(id int64) error {

	var query = "DELETE FROM video_classification WHERE videos_id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deleteVideoClassificationFromVideoGroupId(id int64) error {

	var query = "DELETE FROM video_classification WHERE video_groups_id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deleteVideoFromId(id int64) error {

	if err := deleteVideoClassificationFromVideoId(id); err != nil {
		return err
	}

	if err := deletePermissionsFromGroupId(id); err != nil {
		return err
	}

	var query = "DELETE FROM videos WHERE id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

func deleteVideoGroupFromId(id int64) error {

	if err := deleteVideoClassificationFromVideoGroupId(id); err != nil {
		return err
	}

	if err := deletePermissionsFromVideoGroupId(id); err != nil {
		return err
	}

	var query = "DELETE FROM video_groups WHERE id=?;"
	var params = []interface{}{id}
	return deleteDb(query, params)
}

// Public functions
func UpdateVideoGroup(videoGroup globals.VideoGroup, videos []globals.Video, groups []globals.Group) (globals.VideoGroup, []globals.Video, []globals.Group, error) {

	newVideoGroup, err := updateVideoGroup(videoGroup)
	if err != nil {
		return globals.VideoGroup{}, nil, nil, err
	}

	newVideos, err := updateVideoClassificationFromVideoGroupId(newVideoGroup.Id, videos)
	if err != nil {
		return globals.VideoGroup{}, nil, nil, err
	}

	newGroups, err := updatePermissionsFromVideoGroupId(videoGroup.Id, groups)
	if err != nil {
		return globals.VideoGroup{}, nil, nil, err
	}

	return newVideoGroup, newVideos, newGroups, nil
}

func UpdateVideo(video globals.Video, videoGroups []globals.VideoGroup) (globals.Video, []globals.VideoGroup, error) {

	newVideo, err := updateVideo(video)
	if err != nil {
		return globals.Video{}, nil, err
	}

	newVideoGroups, err := updateVideoClassificationFromVideoId(newVideo.Id, videoGroups)
	if err != nil {
		return globals.Video{}, nil, err
	}
	return newVideo, newVideoGroups, nil
}

func InsertVideoGroup(videoGroup globals.VideoGroup, videos []globals.Video, groups []globals.Group) (globals.VideoGroup, []globals.Video, []globals.Group, error) {

	newVideoGroup, err := insertVideoGroup(videoGroup)
	if err != nil {
		return globals.VideoGroup{}, nil, nil, err
	}

	newVideos, err := updateVideoClassificationFromVideoGroupId(newVideoGroup.Id, videos)
	if err != nil {
		return globals.VideoGroup{}, nil, nil, err
	}

	newGroups, err := updatePermissionsFromVideoGroupId(newVideoGroup.Id, groups)
	return newVideoGroup, newVideos, newGroups, nil
}

func InsertVideo(video globals.Video, videoGroups []globals.VideoGroup) (globals.Video, []globals.VideoGroup, error) {

	newVideo, err := insertVideo(video)
	if err != nil {
		return globals.Video{}, nil, err
	}

	newVideoGroups, err := updateVideoClassificationFromVideoId(newVideo.Id, videoGroups)
	if err != nil {
		return globals.Video{}, nil, err
	}
	return newVideo, newVideoGroups, nil
}

func GetVideoFromId(id int64) (globals.Video, []globals.VideoGroup, error) {

	video, err := getVideoFromId(id)
	if err != nil {
		return globals.Video{}, nil, err
	}
	videoGroups, err := getVideoGroupsFromVideoId(id)
	if err != nil {
		return globals.Video{}, nil, err
	}
	return video, videoGroups, nil
}

func GetVideoFromSlug(slug string) (globals.Video, []globals.VideoGroup, error) {

	video, err := getVideoFromSlug(slug)
	if err != nil {
		return globals.Video{}, nil, err
	}
	videoGroups, err := getVideoGroupsFromVideoId(video.Id)
	if err != nil {
		return globals.Video{}, nil, err
	}
	return video, videoGroups, nil
}

func GetVideoGroupFromId(id int64) (globals.VideoGroup, []globals.Video, []globals.Group, error) {

	videoGroup, err := getVideoGroupFromId(id)
	if err != nil {
		return globals.VideoGroup{}, nil, nil, err
	}
	videos, err := getVideosFromVideoGroupId(id)
	if err != nil {
		return globals.VideoGroup{}, nil, nil, err
	}
	groups, err := getPermissionsFromVideoGroupId(id)
	if err != nil {
		return globals.VideoGroup{}, nil, nil, err
	}
	return videoGroup, videos, groups, nil
}

var GetAllVideos = getAllVideos
var GetAllVideoGroups = getAllVideoGroups
var DeleteVideoFromId = deleteVideoFromId
var DeleteVideoGroupFromId = deleteVideoGroupFromId
