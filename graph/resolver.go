package graph

import (
	"context"
	"fmt"
	"github.com/sidky/photoscout-server/profile"
	"log"
	"strings"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/sidky/photoscout-server/flickr"
)

const UUID = "uuid"

type Resolver struct {
	flickr *flickr.Flickr
}

func NewResolver(flickr *flickr.Flickr) *Resolver {
	return &Resolver{flickr: flickr}
}

func (r *Resolver) Interesting(ctx context.Context, args struct {
	Page *int32
}) *PhotoList {

	response, err := r.flickr.Interesting(args.Page)
	if err != nil {
		log.Print(err)
	}
	photos := make([]*Photo, len(response.Photos.Photos))
	log.Printf("Photos: %d", len(photos))
	for i, photo := range response.Photos.Photos {
		photos[i] = convertPhoto(&photo)
	}

	page := response.Photos.Page.Int()
	var hasNext bool
	if response.Photos.Page == response.Photos.Pages {
		hasNext = false
	} else {
		hasNext = true
	}

	pagination := Pagination{hasNext: hasNext, nextPage: *page + 1}

	return &PhotoList{photos: photos, pagination: &pagination}
}

func (r *Resolver) Search(args struct {
	Query *string
	Bbox  *BoundingBox
	Page  *int32
}) *PhotoList {
	response, err := r.flickr.Search(args.Query, convertBoundingBox(args.Bbox), args.Page)
	if err != nil {
		log.Print(err)
	}
	photos := make([]*Photo, len(response.Photos.Photos))
	log.Printf("Photos: %d", len(photos))
	for i, photo := range response.Photos.Photos {
		photos[i] = convertPhoto(&photo)
	}

	page := response.Photos.Page.Int()
	var hasNext bool
	if response.Photos.Page == response.Photos.Pages {
		hasNext = false
	} else {
		hasNext = true
	}

	pagination := Pagination{hasNext: hasNext, nextPage: *page + 1}

	return &PhotoList{photos: photos, pagination: &pagination}
}

func (r *Resolver) Detail(args struct {
	PhotoId *string
}) *PhotoDetail {
	split := strings.Index(*args.PhotoId, ":")
	source := (*args.PhotoId)[:split]
	log.Printf("Source: %s", source)
	photoId := (*args.PhotoId)[split+1:]
	infoResponse, err := r.flickr.Info(photoId)
	if err != nil {
		log.Print(err)
	}

	info := infoResponse.Photo
	tags := make([]*Tag, len(info.Tags.Tag))
	for index, tag := range info.Tags.Tag {
		gtag := Tag{
			raw:        tag.Raw,
			machineTag: (tag.MachineTag != 0),
		}
		tags[index] = &gtag
	}

	log.Printf("ID from info: %s", info.ID)

	var location *Location
	if loc := info.Location; loc != nil {
		location = &Location{latitude: loc.Latitude.Float(), longitude: loc.Longitude.Float(), accuracy: *loc.Accuracy.Int()}
	}

	exifResponse, err := r.flickr.Exif(photoId)
	photoExif := exifResponse.Photo
	exifs := make([]*Exif, len(photoExif.Tags))
	for index, exif := range photoExif.Tags {
		gexif := Exif{
			tagSpace: exif.TagSpace,
			tag:      exif.Tag,
			label:    exif.Label,
			raw:      exif.Raw.Content,
		}
		exifs[index] = &gexif
	}

	if err != nil {
		log.Print(err)
	}
	owner := PhotoOwner{
		userID:   info.Owner.Nsid,
		name:     info.Owner.RealName,
		userName: info.Owner.UserName,
		location: info.Owner.Location,
	}

	response := PhotoDetail{
		id:          info.ID,
		uploadedAt:  graphql.Time{time.Unix(int64(info.DateUploaded.Int64()), 0)},
		owner:       owner,
		title:       info.Title.Content,
		description: info.Description.Content,
		camera:      &photoExif.Camera,
		tags:        tags,
		exif:        exifs,
		location:    location,
	}
	return &response
}

func (r *Resolver) BookmarkPhoto(ctx context.Context, args struct {
	PhotoId string
}) *OpResult {
	uuid := ctx.Value(UUID).(*string)
	user := profile.User{*uuid}
	err := user.BookmarkPhoto(args.PhotoId)
	if err != nil {
		errMsg := err.Error()
		return &OpResult{success: false, err: &errMsg }
	} else {
		return &OpResult{success: true}
	}
}

func convertPhoto(flickr *flickr.Photo) *Photo {
	var photo Photo

	photo.id = fmt.Sprintf("flickr:%s", flickr.ID)
	photo.ownerName = flickr.OwnerName

	if *flickr.Accuracy.Int() > 0 {
		photo.location = &Location{
			latitude:  flickr.Latitude.Float(),
			longitude: flickr.Longitude.Float(),
			accuracy:  *flickr.Accuracy.Int(),
		}
	}

	urls := make([]*SizedURL, 0)
	urls = append(urls, sizedURL("Thumbnail", flickr.ThumbnailURL, flickr.ThumbnailWidth.Int(), flickr.ThumbnailHeight.Int()))
	urls = append(urls, sizedURL("Small", flickr.SmallURL, flickr.SmallWidth.Int(), flickr.SmallHeight.Int()))
	urls = append(urls, sizedURL("Small320", flickr.Small320URL, flickr.Small320Width.Int(), flickr.Small320Height.Int()))
	urls = append(urls, sizedURL("Square", flickr.SquareURL, flickr.SquareWidth.Int(), flickr.SquareHeight.Int()))
	urls = append(urls, sizedURL("LargeSquare", flickr.LargeSquareURL, flickr.LargeSquareWidth.Int(), flickr.LargeSquareHeight.Int()))
	urls = append(urls, sizedURL("Medium", flickr.MediumURL, flickr.MediumWidth.Int(), flickr.MediumHeight.Int()))
	urls = append(urls, sizedURL("Medium640", flickr.Medium640URL, flickr.Medium640Width.Int(), flickr.Medium640Height.Int()))
	urls = append(urls, sizedURL("Medium800", flickr.Medium800URL, flickr.Medium800Width.Int(), flickr.Medium800Height.Int()))
	urls = append(urls, sizedURL("Large", flickr.LargeURL, flickr.LargeWidth.Int(), flickr.LargeHeight.Int()))
	urls = append(urls, sizedURL("Original", flickr.OriginalURL, flickr.OriginalWidth.Int(), flickr.OriginalHeight.Int()))

	photo.photoURLs = filterNull(urls)

	return &photo
}

func convertBoundingBox(bbox *BoundingBox) *flickr.BoundingBox {
	if bbox != nil {
		return &flickr.BoundingBox{
			MinLongitude: bbox.MinLongitude,
			MinLatitude:  bbox.MinLatitude,
			MaxLongitude: bbox.MaxLongitude,
			MaxLatitude:  bbox.MaxLatitude,
		}
	}
	return nil
}

func sizedURL(size string, url *string, width *int32, height *int32) *SizedURL {
	if url != nil || width != nil || height != nil {
		return &SizedURL{size: size, url: *url, width: *width, height: *height}
	}
	return nil
}

func filterNull(urls []*SizedURL) []*SizedURL {
	filtered := make([]*SizedURL, 0)

	for _, u := range urls {
		if u != nil {
			filtered = append(filtered, u)
		}
	}
	return filtered
}
