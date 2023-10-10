package metadata

type CollectionType string

const (
	TypeDefaultCollection CollectionType = "default_collection"
	TypeCappedCollection  CollectionType = "capped_collection"
	TypeTTLCollection     CollectionType = "ttl_collection"
	// as the time this was written, view creation is defined as collection entity
	TypeViewCollection CollectionType = "view"
)

type CollectionOption string

const (
	CollectionOptionCapped              CollectionOption = "capped"
	CollectionOptionSize                CollectionOption = "size"
	CollectionOptionExpiredAfterSeconds CollectionOption = "expiredAfterSeconds"
)
