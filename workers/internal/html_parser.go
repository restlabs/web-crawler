package internal

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func ParseRawHtmlContent(rawHtml string) error {
	r := strings.NewReader(rawHtml)
	d := xml.NewDecoder(r)

	d.Strict = false
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity

	for {
		tt, err := d.Token()
		switch err {
		case io.EOF:
			return nil
		case nil:
		default:
			return fmt.Errorf("invalid HTML %v. Error %v", tt, err)
		}
	}
}
