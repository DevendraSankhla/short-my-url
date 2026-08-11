package database

type Models struct {
	Urls URlModel
}

func NewModels(db Database) Models {
	return Models{
		Urls: URlModel{db},
	}
}
