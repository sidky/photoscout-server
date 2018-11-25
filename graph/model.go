package graph

import (
	graphql "github.com/graph-gophers/graphql-go"
)

type Pagination struct {
	hasNext  bool
	nextPage int32
}

func (p *Pagination) HasNext() bool {
	return p.hasNext
}

func (p *Pagination) Next() *int32 {
	return &p.nextPage
}

type PhotoList struct {
	photos     []*Photo
	pagination *Pagination
}

func (list *PhotoList) Photos() *[]*Photo {
	return &list.photos
}

func (list *PhotoList) Pagination() *Pagination {
	return list.pagination
}

type Photo struct {
	id        string
	ownerName string
	location  *Location
	photoURLs []*SizedURL
}

func (p *Photo) ID() graphql.ID {
	return graphql.ID(p.id)
}

func (p *Photo) OwnerName() string {
	return p.ownerName
}

func (p *Photo) Location() *Location {
	return p.location
}

func (p *Photo) PhotoURLs() *[]*SizedURL {
	return &p.photoURLs
}

type Location struct {
	latitude  float64
	longitude float64
	accuracy  int32
}

func (l *Location) Latitude() float64 {
	return l.latitude
}

func (l *Location) Longitude() float64 {
	return l.longitude
}

func (l *Location) Accuracy() int32 {
	return l.accuracy
}

type SizedURL struct {
	size   string
	url    string
	width  int32
	height int32
}

func (url *SizedURL) Size() string {
	return url.size
}

func (url *SizedURL) URL() string {
	return url.url
}

func (url *SizedURL) Width() int32 {
	return url.width
}

func (url *SizedURL) Height() int32 {
	return url.height
}

type BoundingBox struct {
	MinLongitude float64
	MinLatitude  float64
	MaxLongitude float64
	MaxLatitude  float64
}

type PhotoDetail struct {
	id          string
	uploadedAt  graphql.Time
	owner       PhotoOwner
	title       string
	description string
	camera      *string
	tags        []*Tag
	exif        []*Exif
}

func (detail *PhotoDetail) ID() string {
	return detail.id
}

func (detail *PhotoDetail) UploadedAt() graphql.Time {
	return detail.uploadedAt
}

func (detail *PhotoDetail) Owner() *PhotoOwner {
	return &detail.owner
}

func (detail *PhotoDetail) Title() string {
	return detail.title
}

func (detail *PhotoDetail) Description() string {
	return detail.description
}

func (detail *PhotoDetail) Tags() *[]*Tag {
	return &detail.tags
}

func (detail *PhotoDetail) Camera() *string {
	return detail.camera
}

func (detail *PhotoDetail) Exif() *[]*Exif {
	return &detail.exif
}

type PhotoOwner struct {
	userID   string
	name     string
	userName string
	location *string
}

func (owner *PhotoOwner) UserID() string {
	return owner.userID
}

func (owner *PhotoOwner) Name() string {
	return owner.name
}

func (owner *PhotoOwner) UserName() string {
	return owner.userName
}

func (owner *PhotoOwner) Location() *string {
	return owner.location
}

type Tag struct {
	raw        string
	machineTag bool
}

func (tag *Tag) Raw() string {
	return tag.raw
}

func (tag *Tag) IsMachineTag() bool {
	return tag.machineTag
}

type Exif struct {
	tagSpace string
	tag      string
	label    string
	raw      string
}

func (exif *Exif) TagSpace() string {
	return exif.tagSpace
}

func (exif *Exif) Tag() string {
	return exif.tag
}

func (exif *Exif) Label() string {
	return exif.label
}

func (exif *Exif) Raw() string {
	return exif.raw
}
