package graph

import (
	"log"

	"github.com/sidky/photoscout-server/flickr"
)

type Resolver struct {
	flickr *flickr.Flickr
}

func NewResolver(flickr *flickr.Flickr) *Resolver {
	return &Resolver{flickr: flickr}
}

func (r *Resolver) Interesting(args struct {
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

func convertPhoto(flickr *flickr.Photo) *Photo {
	var photo Photo

	photo.id = flickr.ID
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
