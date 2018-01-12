package tpdevices

import (
	"reflect"
)

type TPDevice struct {
	Addr string
	Tags []string
	Data map[string]interface{}
}

func (d *TPDevice) initialize() error {
	d.Data = make(map[string]interface{})
	err, tags, data := d.getSystemInfo()
	if err != nil {
		return err
	}
	d.appendData(tags, data)
	return nil
}

func (d *TPDevice) appendData(tags []string, data interface{}) {
	v := reflect.ValueOf(data)
	for _, t := range tags {
		if !d.TagExists(t) {
			d.Tags = append(d.Tags, t)
		}
	}
	for j := 0; j < v.NumField(); j++ {
		field_name := v.Type().Field(j).Name
		if field_name == "ErrorCode" {
			continue
		}
		field_value := v.Field(j).Interface()
		d.Data[field_name] = field_value
	}
}

func (d *TPDevice) TagExists(a string) bool {
	for _, i := range d.Tags {
		if i == a {
			return true
		}
	}
	return false
}

func (d *TPDevice) GetAllData() error {
	v := reflect.TypeOf(d)
	for j := 0; j < v.NumMethod(); j++ {
		method := v.Method(j)
		if method.Name == "Initialize" || method.Name == "GetAllData" || method.Name == "TagExists" {
			continue
		}
		
		result := method.Func.Call([]reflect.Value{reflect.ValueOf(d)})
		err := result[0].Interface()
		tagsR := result[1]
		data := result[2].Interface()

		if err != nil {
			return err.(error)
		}
		tags := make([]string, tagsR.Len())
		for i := 0; i < tagsR.Len(); i++ {
			tags[i] = tagsR.Index(i).String()
		}
		d.appendData(tags, data)
		
	}
	return nil
}