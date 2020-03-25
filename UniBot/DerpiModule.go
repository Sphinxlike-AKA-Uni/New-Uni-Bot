package Uni

import (
	"fmt"
	"time"
	"strings"
	"net/url"
	"encoding/json"
)

type Image struct { // Derpibooru's way to store image info
	ID int `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstSeenAt time.Time `json:"first_seen_at"`
	UploaderID int `json:"uploader_id"`
	FileName string `json:"file_name"`
	Description string `json:"description"`
	Uploader string `json:"uploader"`
	Image string `json:"image"`
	Score int `json:"score"`
	Upvotes int `json:"upvotes"`
	Downvotes int `json:"downvotes"`
	Faves int `json:"faves"`
	CommentCount int `json:"comment_count"`
	Tags string `json:"tags"`
	TagIds []int `json:"tag_ids"`
	Width int `json:"width"`
	Height int `json:"height"`
	AspectRatio float64 `json:"aspect_ratio"`
	OriginalFormat string `json:"original_format"`
	MimeType string `json:"mime_type"`
	Sha512Hash string `json:"sha512_hash"`
	OrigSha512Hash string `json:"orig_sha512_hash"`
	SourceURL string `json:"source_url"`
	Representations struct {
		ThumbTiny string `json:"thumb_tiny"`
		ThumbSmall string `json:"thumb_small"`
		Thumb string `json:"thumb"`
		Small string `json:"small"`
		Medium string `json:"medium"`
		Large string `json:"large"`
		Tall string `json:"tall"`
		Full string `json:"full"`
	} `json:"representations"`
	IsRendered bool `json:"is_rendered"`
	IsOptimized bool `json:"is_optimized"`
	//Interactions []interface{} `json:"interactions"`// assuming it's not important?
}

// Derpi Search Results
type Search struct {
	Search []Image `json:"search"`
	Total int `json:"total"`
	//Interactions []interface{} `json:"interactions"`// assuming it's not important?
}

// For derpi filters
type Filter struct {
	ID int
	Name string
	Description string
	Hidden_Tag_IDs []int
	Spoilered_Tag_IDs []int
	Spoilered_Tags []string
	Hidden_Complex string
	Spoilered_Complex string
	Public bool
	System bool
	User_Count int
	User_ID int
}

// Get a random image with the following tags(and filter if set)
func (Uni *UniBot) SearchOnDerpi(cID, tags string) {
	tags = url.QueryEscape(tags)
	tags = strings.Replace(strings.Replace(tags, " ", "+", -1), "&", "%26", -1)
	fmt.Println(Uni.GetChannelDerpiFilter(cID))
}

// Get the channel's set derpi filter
func (Uni *UniBot) GetChannelDerpiFilter(cID string) (string, error) {
	var fstr string = "157679" // https://www.derpibooru.org/filters/157679 being the default filter
	err := Uni.DBGetFirstVar(fstr, "GetDerpiFilter", cID)
	return fstr, err
}

// Set the channel's derpi filter
func (Uni *UniBot) SetChannelDerpiFilter(gID, cID, filterid string) {
	f, err := Uni.GetDerpiFilter(filterid)
	if f == nil && err == nil { // redirected
		Uni.Respond(cID, "Filter seems to have returned nil, is the filter public?")
		return
	} else if f == nil && err != nil { // genuine error
		Uni.ErrRespond(err, cID, "requesting filter data", map[string]interface{}{"gID": gID, "cID": cID, "err": err, "filterid": filterid})
		return
	}
	// everything is fine
	fstr := ""
	Uni.DBGetFirstVar(fstr, "GetDerpiFilter", cID)
	if fstr == "" { // no such index exists, create index
		_, err = Uni.DBExec("InsertDerpiFilter", gID, cID, filterid)
	} else { // index for channel already exists, update index
		_, err = Uni.DBExec("UpdateDerpiFilter", filterid, cID)
	}
	
	if err != nil {
		Uni.ErrRespond(err, cID, "setting channel derpibooru filter", map[string]interface{}{"gID": gID, "cID": cID, "err": err, "filterid": filterid, "fstr": fstr})
	} else {
		Uni.Respond(cID, fmt.Sprintf("Filter set to ID: %d, %q", f.ID, f.Name))
	}
}

// Get Filter data from derpi
func (Uni *UniBot) GetDerpiFilter(filterid string) (*Filter, error) {
	resp, err := Uni.HTTPRequest("GET", fmt.Sprintf("https://derpibooru.org/filters/%s.json", filterid), map[string]interface{}{"User-Agent": GrabUserAgent(),}, nil)
	if err != nil { return nil, err }
	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 302 { // redirected, likely that the filter didn't exist or isn't public
		return nil, nil
	}
	var f *Filter
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil { return nil, err }
	return f, nil
}


/*
// Grab filter and get filter's data
func GetDerpiFilter(filterid string) (*Filter, error) {
	var f *Filter
	resp, err := Uni.HTTPRequest("GET", link, map[string]interface{}{"User-Agent": GrabUserAgent(true)}, nil)
	if err != nil { return nil, err }
	err = json.NewDecoder(resp.Body).Decode(&f)
	/*
	fmt.Println(filterid)
	body, err := HttpGetRequest(fmt.Sprintf("https://derpibooru.org/filters/%s.json", filterid), "Uni_Derpi_Search")
	if err != nil {
		return nil, err
	}

	var tFil *Filter
	
	json.Unmarshal(body, &tFil)
	return tFil, nil
	* /
}
*/
