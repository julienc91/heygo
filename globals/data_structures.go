package globals

type Video struct {
    Id int64 `json:"id" db:"id"`
    Title string `json:"title" db:"title"`
    Path string `json:"path" db:"path"`
    Slug string `json:"slug" db:"slug"`
    ImdbId string `json:"imdb_id" db:"imdb_id"`
}

type VideoGroup struct {
    Id int64 `json:"id" db:"id"`
    Title string `json:"title" db:"title"`
}

type User struct {
    Id int64 `json:"id" db:"id"`
    Login string `json:"login" db:"login"`
    Password string `json:"password" db:"password"`
    Salt string `json:"salt" db:"salt"`
}

type Group struct {
    Id int64 `json:"id" db:"id"`
    Title string `json:"title" db:"title"`
}

type Invitation struct {
    Id int64 `json:"id" db:"id"`
    Value string `json:"value" db:"value"`
}
