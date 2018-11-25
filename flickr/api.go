package flickr

import (
	"net/url"
	"strconv"

	"github.com/sidky/photoscout-server/rest"
)

const (
	flickrHost     = "https://api.flickr.com/services/rest/"
	keyAPIKey      = "api_key"
	defaultExtras  = "owner_name,o_dims,url_sq,url_t,url_s,url_q,url_m,url_n,url_z,url_c,url_l,url_o,geo"
	keyPage        = "page"
	keyPerPage     = "per_page"
	keyText        = "text"
	keySort        = "sort"
	keyBoundingBox = "bbox"
)

type Flickr struct {
	apiKey   string
	pageSize int
}

func FlickrApi(apiKey string) *Flickr {
	return &Flickr{apiKey: apiKey, pageSize: 20}
}

func (flickr *Flickr) defaultQueryParameters() url.Values {
	values := url.Values{}
	values.Add(keyAPIKey, flickr.apiKey)
	values.Add(keyPerPage, strconv.Itoa(flickr.pageSize))
	values.Add("format", "json")
	values.Add("nojsoncallback", "1")
	return values
}

func (flickr *Flickr) listQuery(method string, page *int32) (*url.URL, error) {
	q := flickr.defaultQueryParameters()
	q.Add("method", method)
	q.Add("extras", defaultExtras)
	if page != nil {
		q.Add(keyPage, strconv.Itoa(int(*page)))
	}
	return rest.BuildURL(flickrHost, q)
}

// TODO: Merge with listQuery()
func (flickr *Flickr) searchQuery(method string, query *string, boundingBox *BoundingBox, page *int32) (*url.URL, error) {
	q := flickr.defaultQueryParameters()
	q.Add("method", method)
	q.Add("extras", defaultExtras)
	q.Add("sort", "interestingness-desc")
	if page != nil {
		q.Add(keyPage, strconv.Itoa(int(*page)))
	}
	if query != nil {
		q.Add(keyText, *query)
	}
	bboxQuery := boundingBox.Query()
	if bboxQuery != nil {
		q.Add(keyBoundingBox, *bboxQuery)
	}
	return rest.BuildURL(flickrHost, q)
}

func (flickr *Flickr) Interesting(page *int32) (*PhotosResponse, error) {
	u, err := flickr.listQuery("flickr.interestingness.getList", page)
	if err != nil {
		return nil, err
	}
	var response PhotosResponse
	if err = rest.Get(u, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (flickr *Flickr) Search(query *string, boundingBox *BoundingBox, page *int32) (*PhotosResponse, error) {
	u, err := flickr.searchQuery("flickr.photos.search", query, boundingBox, page)
	if err != nil {
		return nil, err
	}
	var response PhotosResponse
	if err = rest.Get(u, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (flickr *Flickr) Info(photoId string) (*PhotoInfoResponse, error) {
	q := flickr.defaultQueryParameters()
	q.Add("method", "flickr.photos.getInfo")
	q.Add("photo_id", photoId)
	u, e := rest.BuildURL(flickrHost, q)
	if e != nil {
		return nil, e
	}
	var response PhotoInfoResponse
	if e = rest.Get(u, &response); e != nil {
		return nil, e
	}
	return &response, nil
}

func (flickr *Flickr) Exif(photoId string) (*PhotoExifResponse, error) {
	q := flickr.defaultQueryParameters()
	q.Add("method", "flickr.photos.getExif")
	q.Add("photo_id", photoId)
	u, e := rest.BuildURL(flickrHost, q)
	if e != nil {
		return nil, e
	}
	var response PhotoExifResponse
	if e = rest.Get(u, &response); e != nil {
		return nil, e
	}
	return &response, nil
}
