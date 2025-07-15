<p align="center">
  <img src="https://iili.io/JBuv0g4.png" alt="Logo" height=150>
</p>

# Go-mongr8
> This project is still under initial development. A breaking change might happen anytime. Please, use it while watching for the latest update.

Go-mongr8 is a lightweight migration tool for [MongoDB](https://www.mongodb.com/) written in Go.

The project's philosophy is simplicity and efficiency, achieved by keeping everything at a high level. It simplifies the schema change process in MongoDB, allowing users to focus solely on the latest schema.

## Idea
MongoDB offers flexibility but lacks built-in migration tools. Although it's a schemaless database, Go-mongr8 enables you to manage schema changes directly in your codebase.

**Why use it?**
- Define fields and indexes using a template document, making schema representation clear and consistent.
- Keep an up-to-date schema example in your codebase for easier testing and development.

**Inspired by**:
- **Versioned migrations:** Inspired by the robust and reliable SQL migration management found in frameworks like [Django](https://www.djangoproject.com/).
- **Schema definition:** Emphasizes simplicity, similar to the approach used by [Ent](https://entgo.io/).

## Install

```cli
go get -u github.com/amirkode/go-mongr8@latest
```

For the CLI usage with `go-mongr8` command, you can install globally:
```cli
go install github.com/amirkode/go-mongr8@latest
```

## Usage
Basic operation can be done by using `go-mongr8` CLI command. Complete documentation can be found [here](https://github.com/amirkode/go-mongr8/blob/main/docs/README.md).

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
Complete documentation can be found [here](https://github.com/amirkode/go-mongr8/blob/main/docs/USER_GUIDE.md).

## Features
- Dual-purpose: Use as CLI and as library.
- Lightweight migration tool with one-go CLI.
- Simplified and Descriptive schema declaration for ease of use.
- Covers common MongoDB migration operations.
- And much more coming soon.

#### Main Functionalities
- [x] Init migration folder
- [x] Init schema/collection
- [x] Generate migration files
- [x] Apply migration
- [ ] Rollback migration
#### Wish List
- [ ] Simulate migration
	- preview the planned queries to be executed
- [ ] Consolidate migrations
	- sync the schema with the current database state

For supported MongoDB operations, you can see [here](https://github.com/amirkode/go-mongr8/blob/main/docs/USER_GUIDE.md).

## Limitations
As a disclaimer, This is an unofficial package designed for MongoDB Migration written in golang. As it relies on the MongoDB golang API, please be aware that some functionalities may evolve over time, potentially affecting compatibility with future MongoDB versions.

However, we will continue our efforts to provide support for future updates.

## Contribution
We welcome contributions to **go-mongr8**! If you'd like to contribute, please submit a pull request with your changes. If you find any issues, please report them in the issue tracker.

## License
Copyright (c) 2023-present [Authors](https://github.com/amirkode/go-mongr8/blob/main/AUTHORS) and Contributors. Logo was created by [Bing Chat](https://bing.com).

Go-mongr8 is distributed under the [MIT License](https://opensource.org/license/mit/).
