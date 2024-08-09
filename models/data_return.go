package models

type DataReturn struct {
	Data []interface{}
}

func (d *DataReturn) IsEmpty() bool {
	return d.Data == nil
}
