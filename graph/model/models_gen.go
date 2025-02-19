// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID        string     `json:"id"`
	PostID    string     `json:"postID"`
	Text      string     `json:"text"`
	Author    *User      `json:"author"`
	ReplyTo   *Comment   `json:"replyTo,omitempty"`
	CreatedAt string     `json:"createdAt"`
	Replies   []*Comment `json:"replies,omitempty"`
}

type CreateComment struct {
	Text     string  `json:"text"`
	AuthorID string  `json:"author_id"`
	PostID   string  `json:"post_id"`
	ReplyTo  *string `json:"replyTo,omitempty"`
}

type CreatePost struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	AuthorID      string `json:"author_id"`
	AllowComments bool   `json:"allowComments"`
}

type Mutation struct {
}

type Post struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	Author        *User      `json:"author"`
	CreatedAt     string     `json:"createdAt"`
	AllowComments bool       `json:"allowComments"`
	Comments      []*Comment `json:"comments,omitempty"`
}

type Query struct {
}

type Subscription struct {
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
