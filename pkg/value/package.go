package value

// Product information for order history data integrity.
// Package container dimensions.
// Transient variables.
// todo - use quantity to estimate shipping container dimensions
type Package struct {
	// ids
	Id            string `json:"id,omitempty"`
	ProductId     string `json:"product_id"`
	AddressIdFrom string `json:"address_id_from"`
	AddressIdTo   string `json:"address_id_to"`
	// product data
	ProductName  string `json:"product_name"`
	ProductPrice int64  `json:"product_price"`
	ProductQty   int    `json:"product_qty"`
	ProductImg   string `json:"product_img,omitempty"`
	// usps
	ZipOrigination string `json:"zip_origination"`
	ZipDestination string `json:"zip_destination"`
	// fedex
	ShipperStateCode   string `json:"shipper_state_code"`
	RecipientStateCode string `json:"recipient_state_code"`
	// ups, fedex, usps
	ProductPounds int     `json:"pounds"`
	ProductOunces float32 `json:"ounces"`
	ProductWeight float32 `json:"product_weight"`
	ProductLength int     `json:"product_length"`
	ProductWidth  int     `json:"product_width"`
	ProductHeight int     `json:"product_height"`
	// ups, fedex
	TotalLength int     `json:"length"`
	TotalWidth  int     `json:"width"`
	TotalHeight int     `json:"height"`
	TotalWeight float32 `json:"weight"`
	// vendor data
	VendorName  string `json:"vendor_name"`
	VendorType  string `json:"vendor_type"`
	VendorPrice int64  `json:"vendor_price"`
	TotalPrice  int64  `json:"total_price"`
}
