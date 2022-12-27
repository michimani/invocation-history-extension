package extension

func ExportedExtensionIdentifier(c *Client) string {
	return c.extensionIdentifier
}

func SetExtensionIdentifier(c *Client, id string) {
	c.extensionIdentifier = id
}
