package api_interpreter

// import "github.com/amirkode/go-mongr8/migration/translator/dictionary"

const (
	VarNameDatabase      = "db"
	VarNameCollection    = "collection"
	VarNameCreateOptions = "createOptions"
	VarNameError         = "err"
	VarNameContext       = "ctx"
)

type (
	actionCreateCollection struct {
		Action
	}

	actionCreateIndex struct {
		Action
	}

	actionCreateField struct {
		Action
	}

	actionDropCollection struct {
		Action
	}

	actionDropIndex struct {
		Action
	}

	actionDropField struct {
		Action
	}
)
