# go-mongr8 operations
These are list of supported operation performed by this tool.
### Collection Creation
Collection will be created with several metadata options:
- Capped with size
- Expired after seconds (TTL)

A dummy data will be inserted to maintain the structure integrity, it should cover:
- Fields creation
- Indexes creation

Future supports:
- Cover other collection options
- Schema validation

### Field Creation
This operation expects a single field creation/insertion.
For example:

A field `name` will be added in this document:
```json
{
	"age": 10,
}
```
to:
```json
{
	"age": 10,
    "name": "",
}
```

Note that, in every added field, the value will be set to its default empty format.

This operation also supports field creation on nested object or array.

A field `name` of `age` object will be added in this document:
```json
{
	"age": {
        "value": 10,
    },
}
```
to:
```json
{
	"age": {
        "value": 10,
        "name": "",
    },
}
```

Future supports:
- Maintain defined ordering

### Index Creation
In index cration, no limition of its default operation. it will create the index based given:
- Index Keys
- Index Options

Any type should be supported, such as:
- Single Field Index
- Compound Index
- Text Index
- 2dsphere Index
- Unique Index
- Partial Index
- Collation Index

User's also able to define a raw expression of the index.

### Field Conversion
There are only two conversions allowed:
- Any to string
- Numeric to numeric:
  - double to int64
  - int32 to int64
  - int32 to double
  - int64 to double
Conversions might be supported in the future:
- String to a supported particular type (as long as the string is in the correct format)

#### Simple conversion:
```json
{
	"age": 10,
}
```

to 

```json
{
	"age": "10",
}
```

Query:
```
db.collection.updateMany(
   { },
   [
     {
       $set: {
         age: { $toString: "$age" }
       }
     }
   ]
)
```

#### Basic nested conversion
```json
{
	"other": {
		"age": 10
	}
}
```

to 

```json
{
	"other": {
		"age": "10"
	}
}
```

Query:
```
db.collection.updateMany(
   { },
   [
     {
       $set: {
         "other.age": { $toString: "$other.age" }
       }
     }
   ]
)
```

#### Conversion on an array of object:
```json
{
	"other": {
		"ages": [
			{
				"local": 10,
				"international": 11,
			},
			{
				"local": 11,
				"international": 11,
			},
		]
	}
}
```

to 

```json
{
	"other": {
		"ages": [
			{
				"local": "10",
				"international": "11",
			},
			{
				"local": "11",
				"international": "11",
			},
		]
	}
}
```

Query:
```
db.collection.updateMany(
  {},
  [{ $set: {
    "other.ages": {
      $map: {
        input: "$other.ages", 
        as: "age",
        in: {
          local: { $toString: "$$age.local" },
          international: { $toString: "$$age.international" }
        }
      }
    }
  }}]
)
```

#### Conversion on an array of scalar:
```json
{
	"other": {
		"scores": [
			10,
			11,
		]
	}
}
```

to 

```json
{
	"other": {
		"scores": [
			"10",
			"11",
		]
	}
}
```

Query:
```
db.collection.updateMany(
  {},
  [{
    $set: {
      "other.scores": {
        $map: {
          input: "$other.scores",
          as: "score",
          in: { $toString: "$$score" }
        }
      }
    }
  }]
)
```

#### Coversion on direct Nested Array
```json
{
	"other": {
		"scores": [
			[1, 2, 3, 4],
			[5, 6, 7, 8]
		]
	}
}
```

to

```json
{
	"other": {
		"scores": [
			["1", "2", "3", "4"],
			["5", "6", "7", "8"]
		]
	}
}
```

Query:
```
db.users.updateMany({},
[
  {
    $set: {
      "other.scores": {
        $map: {
          input: "$other.scores",
          as: "outer",  
          in: {
            $map: {
              input: "$$outer",
              as: "inner",
              in: {$toString: "$$inner"}
            }
          }
        }
      }
    }
  }
])
```

#### Conversion on direct Nested 3 levels array
```json
{
	"other": {
		"scores": [
			[[1, 2], [3, 4]],
			[[5, 6], [7, 8]]
		]
	}
}
```

to

```json
{
	"other": {
		"scores": [
			[["1", "2"], ["3", "4"]],
			[["5", "6"], ["7", "8"]]
		]
	}
}
```

Query:
```
db.collection.updateMany({},
[
  {
    $set: {
      "other.scores": {
        $map: {
          input: "$other.scores",
          as: "outer",
          in: {
            $map: {
              input: "$$outer",
              as: "middle",
              in: {
                $map: {
                  input: "$$middle",
                  as: "inner",
                  in: {$toString: "$$inner"}
                }
              }
            }
          }
        }  
      }
    }
  }
])
```


#### Conversion on Nested Array inside object
```json
{
	"other": {
		"scores": [
			{
				"taken_at": Date("2023-01-01"),
				"values": [1, 2, 3, 4]
			}
		]
	}
}
```

to 

```json
{
	"other": {
		"scores": [
			{
				"taken_at": "2023-01-01",
				"values": ["1", "2", "3", "4"]
			}
		]
	}
}
```

Query:
```
db.users.updateMany({}, 
[
  {
    $set: {
      "other.scores": {
        $map: {
          input: "$other.scores",
          as: "score", 
          in: {
            taken_at: { $dateToString: { format: "%Y-%m-%d", date: "$$score.taken_at"} },
            values: {
              $map: {
                input: "$$score.values",
                as: "val",
                in: { $toString: "$$val" }
              }
            }
          }
        }
      }
    }
  }
])
```

#### Conversion on nested array of object
```json
{
	"other": {
		"scores": [
			{
				"values": [
					{
						"previous": 10,
						"current": 11
					},
					{
						"previous": 10,
						"current": 12
					},
				]
			}
		]
	}
}
```

to 

```json
{
	"other": {
		"scores": [
			{
				"values": [
					{
						"previous": "10",
						"current": "11"
					},
					{
						"previous": "10",
						"current": "12"
					},
				]
			}
		]
	}
}
```

Query:
```
db.users.updateMany({}, 
[
  {
    $set: {
      "other.scores": {
        $map: {
          input: "$other.scores",
          as: "score",
          in: {
            values: {
              $map: {
                input: "$$score.values",
                as: "val",
                in: {
                  previous: {$toString: "$$val.previous"},
                  current: {$toString: "$$val.current"}
                }
              }
            }  
          }
        }
      }
    }
  }
])
```

#### Conversion on nested array of nested object
```json
{
	"other": {
		"ages": [
			{
				"local": {
					"value": {
						"value1": 1
					}
				},
				"international": {
					"value": 11
				},
			}
		]
	}
}
```

to

```json
{
	"other": {
		"ages": [
			{
				"local": {
					"value": {
						"value1": "1"
					}
				},
				"international": {
					"value": "11"
				},
			}
		]
	}
}
```

Query:
```
db.users.updateMany({},
[
  {
    $set: {
      "other.ages": {
        $map: {
          input: "$other.ages",
          as: "age", 
          in: {
            local: {
              value: {
                value1: {$toString: "$$age.local.value.value1"}
              }
            },
            international: {
              value: {$toString: "$$age.international.value"} 
            }
          }
        }  
      }
    }
  }
])
```