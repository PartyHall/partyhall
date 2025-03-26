package models

import "fmt"

type BackdropAlbum struct {
	Id      int64  `db:"id" json:"id"`
	NexusId int64  `db:"nexus_id" json:"nexus_id"`
	Name    string `db:"name" json:"name"`
	Author  string `db:"author" json:"author"`
	Version int64  `db:"version" json:"version"`

	Backdrops []Backdrop `db:"-" json:"backdrops"`
}

func (ba BackdropAlbum) String() string {
	return fmt.Sprintf("BackdropAlbum[%v by %v (v%v) - Nexus ID: %v]", ba.Name, ba.Author, ba.Version, ba.NexusId)
}

type Backdrop struct {
	Id       int64  `db:"id" json:"id"`
	AlbumId  int64  `db:"album_id" json:"album_id"`
	NexusId  int64  `db:"nexus_id" json:"nexus_id"`
	Title    string `db:"title" json:"title"`
	Filename string `db:"filename" json:"filename"`
}

func (b Backdrop) String() string {
	return fmt.Sprintf("Backdrop[%v (Nexus ID: %v)]", b.Title, b.NexusId)
}
