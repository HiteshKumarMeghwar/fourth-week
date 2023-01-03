package forms

import (
	"net/url"
)

// Form creates a custom from struct, embeds url.Values object
type Form struct {
	url.Values
	Error errors
}

// New initializes a form struct
/* func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
} */

// Has checks form field is in post and no empty
/* func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}
*/
