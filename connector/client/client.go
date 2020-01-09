package client

type Client interface {
	Post(url url.URL, activity schema.Activity) error
}
