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
- **Capped**, declaration: `[base metadata].Capped([capped size])`.
- **TTL**, declaration: `[base metadata].TTL([expired after seconds])`.

Example:
```go
func (Users) Collection() collection.Metadata {
	return metadata.InitMetadata("users").Capped(1000).TTL(60)
}
```
### Field
import: `github.com/amirkode/collection/field`

Available Fields:
- **String**, declaration: `field.StringField("[field name]")`.
- **Int32**, declaration: `field.Int32Field("[field name]")`.
- **Int64**, declaration: `field.Int64Field("[field name]")`.
- **Double**, declaration: `field.DoubleField("[field name]")`.
- **Boolean**, declaration: `field.BooleanField("[field name]")`.
- **Array**, declaration: `field.ArrayField("[field name]", [child field as array type])`.

  for example:
  ```go
   field.ArrayField("login_history", 
     field.StringField(""),
   )
  ```
  Note that the child field does not require field name.
- **Object**, declaration: `field.ObjectField("[field name]", [children]) `
  for example:
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
- **Timestamp**, declaration: `field.TimestampField("[field name]")`.
- **Geo JSON Point**, declaration: `field.GeoJSONPointField("[field name]")`.
- **Geo JSON Line String**, declaration: `field.GeoJSONLineStringField("[field name]")`.
- **Geo JSON Polygon Single Ring**, declaration: `field.GeoJSONPolygonSingleRingField("[field name]")`.
- **Geo JSON Polygon Multiple Ring**, declaration: `field.GeoJSONPolygonMultipleRingField("[field name]")`.
- **Geo JSON Multi Point**, declaration: `field.GeoJSONMultiPointField("[field name]")`.
- **Geo JSON Multi Line String**, declaration: `field.GeoJSONMultiLineStringField("[field name]")`.
- **Geo JSON Geometry**, declaration: `field.GeoJSONGeometryCollectionField("[field name]")`.
- **Legacy Coordinate Array**, declaration: `field.LegacyCoordinateArrayField("[field name]")`.

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

Available Indexes:
- **Single Field Index**, declaration: `index.SingleFieldIndex(index.Field("field name", [value]))`
- **Compound Index**, declaration: `index.CompoundIndex([fields/keys definition])`

  for example:
  ```go
  index.CompoundIndex(
      index.Field("name", -1),
	  index.Field("age", 1),
  ),
  ```
- **Text Index**, declaration: `index.TextIndex(index.Field("field name", [value]))`
- **Geospatial 2dsphere Index**, declaration: `index.Geospatial2dsphereIndex(index.Field("field name", [value does not matter]))`
- **Unique Index**, declaration: `index.UniqueIndex(index.Field("field name", [value]))`
- Coming soon for Partial and Collation