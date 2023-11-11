# Go-mongr8 User Guide
Go-mongr8 provides various simplicity and eficiency for managing MongoDB schema migration so that you can simply declare high-level schema definition. We provide most APIs to the schema.

Here's the schema definition example:
```go
func (Users) Collection() collection.Metadata {
	return metadata.InitMetadata("users")
}

func (Users) Fields() []collection.Field {
	return []collection.Field{
		field.StringField("name"),
		field.Int32Field("age"),
	}
}

func (Users) Indexes() []collection.Index {
	return []collection.Index{
		index.SingleFieldIndex(index.Field("name", 1)),
		index.CompoundIndex(
			index.Field("name", -1),
			index.Field("age", 1),
		),
	}
}

```

As shown above, you can declare basic schema definition such as collection **Metadata**, **Fields**, and **Indexes**.

### Supported MongoDB Operations
- [x] Create collection with options:
  - [x] Capped
  - [x] Expiration (TTL)
  - others coming soon
- [x] Create field (in any depth):
  - [x] String
  - [x] Int32
  - [x] Int64
  - [x] Double
  - [x] Boolean
  - [x] Array
  - [x] Date (Timestamp)
  - [x] Geo JSON Point
  - [x] Geo JSON Line String
  - [x] Geo JSON Polygon Single Ring
  - [x] Geo JSON Polygon Multiple Ring
  - [x] Geo JSON Multi Point
  - [x] Geo JSON Multi Line String
  - [x] Geo JSON Multi Polygon
  - [x] Geo JSON Geometry Collection
  - [x] Legacy Coordinate Array
- [x] Create index:
  - [x] Single Field
  - [x] Compound
  - [x] Text
  - [x] Geospatial 2dsphere
  - [x] Unique
  - [x] Partial
  - [x] Collation
  - [x] Raw Expression
- [x] Field type conversion (in any depth):
    - [x] Number to number
    - [x] Any to string
    - [ ] String to any (in usecase validation)
- [x] Drop Collection
- [x] Drop Field (in any depth)
- [x] Drop Index
- [ ] Auto Apply Schema Validation (soon)

## Getting Started
Please ensure that you have already initiated the `go-mongr8` in your project. Complete documentation can be found [here]().

### Define Schema
Now, navigate to your project and create a collection schema template, let's taken an example of `users` collection.

```sh
> go-mongr8 create-collection users
```

After that, you can open the template in `mongr8/collection/users.go` (note that all schema template will be in snake case naming).

```go
func (Users) Collection() collection.Metadata {
    // you may add additional options here
	return metadata.InitMetadata("users")
}

func (Users) Fields() []collection.Field {
	return []collection.Field{
		// define the fields here
	}
}

func (Users) Indexes() []collection.Index {
	return []collection.Index{
		// define the indexes here
	}
}

```

You define the schema with available `go-mongr8` APIs. Here some explanation of the components:
- **Metadata** contains basic informations of a collection such as name, and options. You can see available APIs [here]().
- **Fields** are the list of document outer fields. We supports most primitive types as well as complex types such as object, array, and some Geo JSON fields. You can see available APIs [here]().
- **Indexes** are the definition of collection indexes. We supports any supported MongoDB index type (at least as the project released). You can see available APIs [here]().

### Generate Migration Files
After defining schema, you need to generate migration file. It will generate a migration with particular version (by default it uses the timestamp when the command executed). So, it makes sure that every time changes are made will be grouped in different versions.

You can start generate by executing:
```sh
> go-mongr8 generate-migration
```
You may add description to the migration:
```sh
> go-mongr8 generate-migration --desc "my migration description"
```

This will produce a new migration file in `mongr8/migration`. Please do not change anything to the generated code.

The migration file should contain relevant actions based on the changes of the latest schema. Here are possible actions:
- Create Collection
- Create Field
- Create Index
- Convert Field
- Drop Collection
- Drop Field
- Drop Index

### Apply Migrations
After having migration files ready, please make sure required configurations are set as stated [here]().

You can apply the migrations by executing:
```sh
> go-mongr8 apply-migration
```
If you have a replicaset server, you can apply the migrations whithin a transaction session by adding this flag:
```sh
> go-mongr8 apply-migration --use-transaction
```
You can check whether migrations are applied by simply checking the history on `mongr8_migration_history` collection.

It's worth nothing, Go-mongr8 always maintains a **dummy document** in each collection.

## Go-mongr8 APIs
### Metadata
import: `github.com/amirkode/collection/metadata`

Basic declaration: 
```go
metadata.InitMetadata("[collection name]")
```

All Options:
- **Capped**

	Declaration:
	```go
	[base metadata].Capped([capped size])
	```
- **TTL**

	Declaration: 
	```go
	[base metadata].TTL([expired after seconds])
	```

- Example:
	```go
	func (Users) Collection() collection.Metadata {
		return metadata.InitMetadata("users").Capped(1000).TTL(60)
	}
	```
### Field
import: `github.com/amirkode/collection/field`

Available Fields:
- **String**
	
	Declaration:
	```go
	field.StringField("[field name]")
	```
- **Int32**
	
	Declaration:
	```go
	field.Int32Field("[field name]")
	```
- **Int64**
	
	Declaration:
	```go
	field.Int64Field("[field name]")
	```
- **Double**
	
	Declaration:
	```go
	field.DoubleField("[field name]")
	```
- **Boolean**
	
	Declaration:
	```go
	field.BooleanField("[field name]")
	```
- **Array**
	
	Declaration:
	```go
	field.ArrayField("[field name]", [child field as array type])
	```

	For example:
	```go
	field.ArrayField("login_history", 
		field.StringField(""),
	)
	```
	Note that the child field does not require field name.
- **Object**
	
	Declaration:
	```go
	field.ObjectField("[field name]", [children])
	```
	For example:
	```go
	field.ObjectField("other_info",
		field.StringField("full_address"),
		field.StringField("city"),
		field.StringField("country"),
		field.ArrayField("scores",
			field.Int32Field(""),
	),
	)
	```  
- **Timestamp**
	
	Declaration:
	```go
	field.TimestampField("[field name]")
	```
- **Geo JSON Point**
	
	Declaration:
	```go
	field.GeoJSONPointField("[field name]")
	```
- **Geo JSON Line String**
	
	Declaration:
	```go
	field.GeoJSONLineStringField("[field name]")
	```
- **Geo JSON Polygon Single Ring**
	
	Declaration:
	```go
	field.GeoJSONPolygonSingleRingField("[field name]")
	```
- **Geo JSON Polygon Multiple Ring**
	
	Declaration:
	```go
	field.GeoJSONPolygonMultipleRingField("[field name]")
	```
- **Geo JSON Multi Point**
	
	Declaration:
	```go
	field.GeoJSONMultiPointField("[field name]")
	```
- **Geo JSON Multi Line String**
	
	Declaration:
	```go
	field.GeoJSONMultiLineStringField("[field name]")
	```
- **Geo JSON Geometry**

	Declaration:
	```go
	field.GeoJSONGeometryCollectionField("[field name]")
	```
- **Legacy Coordinate Array**
	
	Declaration:
	```go
	field.LegacyCoordinateArrayField("[field name]")
	```

### Index
import: `github.com/amirkode/collection/index`

Generally the declaration of field is in this format:
```
AnIndexAPI([field/key definition])
```
Here's the field definition API:
```go
index.Field("field name", "[value]")
// for example:
index.Field("name", 1)
```
Sometimes, we need to create an index for nested field. You can use this chaining method:
To apply this index key:
```json
{
	"other.info": 1,
}
```
```go
// you can basically declare
index.Field("other.info.value", 1)
// or you can chain this method many times to make it more idiomatic
index.Field("other", 1).NestedField("info").NestedField("info")
```

#### Available Indexes
- **Single Field Index**
	
	Declaration:
	```go
	index.SingleFieldIndex(index.Field("field name", [value]))
	```
	To create a single field index, you don't necessarily declare using this function. You may also declare using Compound Index or Raw Index API.

- **Compound Index**
	
	Declaration:
	```go
	index.CompoundIndex([fields/keys definition])
	```
	For example:
	```go
	index.CompoundIndex(
		index.Field("name", -1),
		index.Field("age", 1),
	),
	```
- **Text Index**
	
	Declaration:
	```go
	index.TextIndex(index.Field("field name"))
	```
- **Geospatial 2dsphere Index**
	
	Declaration:
	```go
	index.Geospatial2dsphereIndex(index.Field("field name"))
	```
- **Hashed Index**
	
	Declaration:
	```go
	index.HashedIndex(index.Field("field name"))
	```

#### Option
You may also set several options on the index. Our API also provide most options.

The option can added by method chanining, here's the format:
```go
// [index declaration].[option method]
index.SingleFieldIndex(index.Field("name", 1)).AsSparse()
```
Available Options:
- **Sparse**
	
	Declaration:
	```go
	[index].AsSparse()
	```
	It will add this to the option map:
	```json
	{
		"sparse": true
	}
	```
- **Background**

	Declaration:
	```go
	[index].AsBackground()
	```
	It will add this to the option map:
	```json
	{
		"background": true
	}
	```
- **Hidden**

	Declaration:
	```go
	[index].AsHidden()
	```
	It will add this to the option map:
	```json
	{
		"hidden": true
	}
	```
- **Partial Expression**
	
	Declaration:
	```go
	[index].SetPartialExpression([a map of interface])
	```
	For example:
	```go
	[index].SetPartialExpression(map[string]interface{}{
		"score": map[string]interface{}{
			"$gt": 70,
		},
	})
	```
	It will add this to the option map:
	```json
	{
		"partialFilterExpression": [argument as a map of interface]
	}
	```
- **TTL**
	
	Declaration:
	```go
	[index].SetTTL([an integer])
	```
	For example:
	```go
	[index].SetTTL(60)
	```
	It will add this to the option map:
	```json
	{
		"expireAfterSeconds": [an integer value]
	}
	```
- **Collation**
	
	Declaration:
	```go
	[index].SetCollation([a map of interface])
	```
	For example:
	```go
	[index].SetCollation(map[string]interface{}{
		"locale": "en_US",
	})
	```
	It will add this to the option map:
	```json
	{
		"collation": [argument as a map of interface]
	}
	```
- **Index Name**

	By default, our index API will generated an index name based on an arrangement of index keys and options. But, We also provide custom index name setting.
	Declaration:
	```go
	[index].SetCustomIndexName("index name goes here")
	```

The provided APIs might not be enough for a few cases. You can also declare a raw expression of an index:
- **Raw Index Expression**
	
	Declaration:
	```go
	[index].RawIndex([a map of interface for index keys], [a pointer map of interface for index options])
	```
	For example:
	```go
	// without option
	[index].RawIndex(
		map[string]interface{}{
			"name": 1,
			"age": -1,
		},
		nil,
	)
	
	// with option
	[index].RawIndex(
		map[string]interface{}{
			"name": 1,
			"age": -1,
		},
		&map[string]interface{}{
			"unique": true,
			"background": true,
		},
	)
	```