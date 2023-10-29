package metadata

type CollectionType string

const (
	TypeDefaultCollection CollectionType = "TypeDefaultCollection"
	// TypeCappedCollection  CollectionType = "capped_collection"
	// TypeTTLCollection     CollectionType = "ttl_collection"
	// as the time this was written, view creation is defined as collection entity
	TypeViewCollection CollectionType = "TypeViewCollection"
)

type CollectionOption string

const (
	CollectionOptionCapped              CollectionOption = "capped"
	CollectionOptionCappedSize          CollectionOption = "size"
	CollectionOptionExpiredAfterSeconds CollectionOption = "expiredAfterSeconds"
)
