package drawing

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type DrawPart struct {
	Tag         string `json:"_tag"`
	Coordinates []int  `json:"coordinates"`
	Times       []int  `json:"times"`
	Controls    []int  `json:"controls",omitempty`
	LineWidth   int    `json:"lineWidth"`
	StrokeStyle Rgba   `json:"strokeStyle"`
	FillStyle   Rgba   `json:"fillStyle",omitempty`
}

var DefaultDrawPart DrawPart = DrawPart{
	Tag:         "",
	Coordinates: nil,
	Times:       nil,
	Controls:    nil,
	LineWidth:   1,
	StrokeStyle: rgba(0, 0, 0, 1),
	FillStyle:   rgba(0, 0, 0, 0),
}

func (d DrawPart) ToJSON() ([]byte, error) {
	resultObj := make(map[string]interface{})
	resultObj["_tag"] = d.Tag
	resultObj["coordinates"] = d.Coordinates
	resultObj["times"] = d.Times
	if d.Controls != nil && len(d.Controls) > 0 {
		resultObj["controls"] = d.Controls
	}
	if d.StrokeStyle != DefaultDrawPart.StrokeStyle {
		resultObj["strokeStyle"] = d.StrokeStyle
	}
	if d.FillStyle != DefaultDrawPart.FillStyle {
		resultObj["fillStyle"] = d.FillStyle
	}
	return json.Marshal(resultObj)
}

func (d *DrawPart) FromJSON(data []byte) error {
	*d = DefaultDrawPart
	return json.Unmarshal(data, d)
}

func rgba(r, g, b uint8, a float32) Rgba {
	return Rgba{[3]uint8{r, g, b}, a}
}

type Rgba struct {
	Rgb   [3]uint8
	Alpha float32
}

func (r Rgba) MarshalJSON() ([]byte, error) {
	asStr := fmt.Sprintf(
		"\"rgba(%d,%d,%d,%.2f)\"",
		r.Rgb[0],
		r.Rgb[1],
		r.Rgb[2],
		r.Alpha)
	return []byte(asStr), nil
}

func (c *Rgba) UnmarshalJSON(bs []byte) error {
	str := string(bs)
	err := errors.New("Bad RGBA format")
	if len(str) == 0 ||
		(strings.Index(str, "\"rgba(") != 0 &&
			strings.LastIndex(str, ")\"") != len(str)-2) {
		return err
	}
	numbers := strings.Split(strings.Trim(str, "\"()rgba"), ",")
	if len(numbers) != 4 {
		return err
	}

	r, rErr := strconv.Atoi(numbers[0])
	g, gErr := strconv.Atoi(numbers[1])
	b, bErr := strconv.Atoi(numbers[2])
	a, aErr := strconv.ParseFloat(numbers[3], 32)
	if rErr == nil &&
		gErr == nil &&
		bErr == nil &&
		aErr == nil {
		c.Rgb = [3]uint8{uint8(r), uint8(g), uint8(b)}
		c.Alpha = float32(math.Min(1, math.Max(0, a)))
		return nil
	}
	return err
}

var validTags map[string]bool = map[string]bool{
	"ble._2d.StrokeReplay":   true,
	"ble._2d.EraseReplay":    true,
	"ble._2d.PolylineReplay": true,
}

func isTagValid(tag string) bool {
	return validTags[tag]
}
