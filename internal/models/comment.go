package models

type Comment struct {
    ID     int    `json:"id"`
    PostID int    `json:"post_id"`
    Author string `json:"author"`
    Text   string `json:"text"`
}
