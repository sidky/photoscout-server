package flickr

import (
	"fmt"
	"net/url"

	"github.com/sidky/photoscout-server/common"
)

// PhotosResponse : Response of search and interesting photo list, returned by Flickr
type PhotosResponse struct {
	Photos PhotoList `json:"photos"`
	Stat   string    `json:"stat"`
}

// PhotoList : List of photos returned by Flickr
type PhotoList struct {
	Page    common.FlexInt `json:"page"`
	Pages   common.FlexInt `json:"pages"`
	PerPage common.FlexInt `json:"perpage"`
	Total   common.FlexInt `json:"total"`
	Photos  []Photo        `json:"photo"`
}

// Photo : A photo object as defined by Flickr
type Photo struct {
	ID        string           `json:"id"`
	OwnerID   string           `json:"owner"`
	Title     string           `json:"title"`
	OwnerName string           `json:"ownername"`
	Latitude  common.FlexFloat `json:"latitude"`
	Longitude common.FlexFloat `json:"longitude"`
	Accuracy  common.FlexInt   `json:"accuracy"`

	SquareURL    *string         `json:"url_sq"`
	SquareHeight *common.FlexInt `json:"height_sq"`
	SquareWidth  *common.FlexInt `json:"width_sq"`

	ThumbnailURL    *string         `json:"url_t"`
	ThumbnailHeight *common.FlexInt `json:"height_t"`
	ThumbnailWidth  *common.FlexInt `json:"width_t"`

	SmallURL    *string         `json:"url_s"`
	SmallHeight *common.FlexInt `json:"height_s"`
	SmallWidth  *common.FlexInt `json:"width_s"`

	LargeSquareURL    *string         `json:"url_q"`
	LargeSquareHeight *common.FlexInt `json:"height_q"`
	LargeSquareWidth  *common.FlexInt `json:"width_q"`

	MediumURL    *string         `json:"url_m"`
	MediumHeight *common.FlexInt `json:"height_m"`
	MediumWidth  *common.FlexInt `json:"width_m"`

	Small320URL    *string         `json:"url_n"`
	Small320Height *common.FlexInt `json:"height_n"`
	Small320Width  *common.FlexInt `json:"width_n"`

	Medium640URL    *string         `json:"url_z"`
	Medium640Height *common.FlexInt `json:"height_z"`
	Medium640Width  *common.FlexInt `json:"width_z"`

	Medium800URL    *string         `json:"url_c"`
	Medium800Height *common.FlexInt `json:"height_c"`
	Medium800Width  *common.FlexInt `json:"width_c"`

	LargeURL    *string         `json:"url_l"`
	LargeHeight *common.FlexInt `json:"height_l"`
	LargeWidth  *common.FlexInt `json:"width_l"`

	OriginalURL    *string         `json:"url_o"`
	OriginalHeight *common.FlexInt `json:"height_o"`
	OriginalWidth  *common.FlexInt `json:"width_o"`
}

type BoundingBox struct {
	MinLongitude float64
	MinLatitude  float64
	MaxLongitude float64
	MaxLatitude  float64
}

func (b *BoundingBox) Query() *string {
	if b != nil {
		q := url.QueryEscape(fmt.Sprintf("%f,%f,%f,%f", b.MinLongitude, b.MinLatitude, b.MaxLongitude, b.MaxLatitude))
		return &q
	}
	return nil
}
