<p align="center">
  <img src="https://iili.io/JBuv0g4.png" alt="Logo" height=150>
</p>

# Go-mongr8
> This project is still under initial development. A breaking change might happen anytime. Please, use it while watching for the latest update.

Go-mongr8 is a lightweight migration tool for [MongoDB](https://www.mongodb.com/) written in Go. This was inspired by SQL migration management in most application frameworks (i.e [Django](https://www.djangoproject.com/)).

The project's philosophy is simplicity and efficiency, achieved by keeping everything at a high level. It simplifies the schema change process in MongoDB, allowing users to focus solely on the latest schema.

## Install

```cli
go get -u github.com/amirkode/go-mongr8@latest
```

For the CLI usage with `go-mongr8` command, you can install globally:
```cli
go install github.com/amirkode/go-mongr8@latest
```

## Usage
Basic operation can be done by using `go-mongr8` CLI command. Complete documentation can be found [here](https://github.com/amirkode/go-mongr8/blob/main/doc/README.md).

Users can easily define their schema with provided APis in this package. Here's the example of simple schema definition:

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
Complete documentation can be found [here](https://github.com/amirkode/go-mongr8/blob/main/doc/USER_GUIDE.md).

## Features
- Dual-purpose: Use as CLI and as library.
- Lightweight migration tool with one-go CLI.
- Simplified and Descriptive schema declaration for ease of use.
- Covers common MongoDB migration operations.
- And much more coming soon.

#### Functionalities
- [x] Init migration folder
- [x] Init schema/collection
- [x] Generate migration file
- [x] Apply migration
- [ ] Consolidate migration
- [ ] Rollback migration

For supported MongoDB operations, you can see [here](https://github.com/amirkode/go-mongr8/blob/main/doc/USER_GUIDE.md).

## Limitations
As a disclaimer, This is an unofficial package designed for MongoDB Migration written in golang. As it relies on the MongoDB golang API, please be aware that some functionalities may evolve over time, potentially affecting compatibility with future MongoDB versions.

However, we will continue our efforts to provide support for future updates.

## Contribution
Coming soon
## License
Copyright (c) 2023-present-present [Authors](https://github.com/amirkode/go-mongr8/blob/main/AUTHORS) and Contributors. Logo was created by [Bing Chat](https://bing.com).

Go-mongr8 is distributed under the [MIT License](https://opensource.org/license/mit/).
