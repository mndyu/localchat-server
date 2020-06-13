package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

type Jsonmap map[string]interface{}

func getJSONMap(value interface{}) {
	// handler.structToJsonmap
}

func convertJSONTypes(value Jsonmap) {
	//testutils.migration
}

func GetFieldByJSONField() {

}

func InterfaceToJsonmap(a interface{}) Jsonmap {
	switch v := a.(type) {
	case Jsonmap:
		return v
	case map[string]interface{}:
		return Jsonmap(v)
	}
	return nil
}

func StructToJsonmap(out *Jsonmap, a interface{}) error {
	EachSchemaField(a, func(jsonFieldName string, val reflect.Value, field reflect.StructField) {
		(*out)[jsonFieldName] = val.Interface()
	})
	return nil
}
func AssignMapToStruct(out interface{}, a Jsonmap) error {
	EachSchemaField(out, func(jsonFieldName string, val reflect.Value, field reflect.StructField) {
		v := val //reflect.ValueOf(i)
		if !v.IsValid() || !v.CanSet() {
			// fmt.Println("cant set:", uf.Name)
			return
		}
		fieldValue := a[jsonFieldName]
		if fieldValue == nil {
			return
		}
		fv := reflect.ValueOf(fieldValue)
		err := AssignJSONValue(v, fv)
		if err != nil {
			t1 := fv.Type().String()
			t2 := v.Type().String()
			fmt.Println(err, t1, t2)
			return
			// return fmt.Errorf("field %s of type %s is not assignable to type %s", uf.Name, ufv.Type().Name(), fv.Type().Name())
		}
	})
	return nil
}

func FilterJsonmapWithStruct(jm Jsonmap, i interface{}) Jsonmap {
	newMap := Jsonmap{}
	EachSchemaField(i, func(jsonFieldName string, val reflect.Value, field reflect.StructField) {
		if jm[jsonFieldName] != nil {
			newMap[jsonFieldName] = jm[jsonFieldName]
		}
	})
	return newMap
}

func isNormalStruct(v reflect.Value) bool {
	k := v.Type().Kind()
	if k != reflect.Struct {
		return false
	}
	i := v.Interface()
	if _, ok := i.(time.Time); ok {
		return false
	}
	if _, ok := i.(decimal.Decimal); ok {
		return false
	}
	return true
}

func EachSchemaField(a interface{}, f func(jsonFieldName string, val reflect.Value, field reflect.StructField)) {
	var uv reflect.Value
	if v, ok := a.(reflect.Value); ok {
		uv = v
	} else {
		uv = reflect.ValueOf(a)
	}
	if uv.Type().Kind() == reflect.Ptr || uv.Type().Kind() == reflect.Slice || uv.Type().Kind() == reflect.Array {
		uv = uv.Elem()
	}
	ut := uv.Type()

	for i := 0; i < ut.NumField(); i++ {
		uf := ut.Field(i)
		ufv := uv.Field(i)
		if uf.Anonymous || isNormalStruct(ufv) {
			EachSchemaField(ufv, f)
			continue
		}
		if strings.ToLower(uf.Name) == uf.Name {
			// ignore private fields
			continue
		}
		fieldName := uf.Tag.Get("json")
		if fieldName == "" {
			fieldName = uf.Name
		}
		f(fieldName, ufv, uf)
	}
}

func ParseGormTagSetting(tags reflect.StructTag) map[string]string {
	setting := map[string]string{}
	for _, str := range []string{tags.Get("sql"), tags.Get("gorm")} {
		if str == "" {
			continue
		}
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}

func GetQueryStructValues(a interface{}) []reflect.Value {
	fields := []reflect.Value{}
	EachQueryStructField(a, func(jsonFieldName string, val reflect.Value, field reflect.StructField) {
		fields = append(fields, val)
	})
	return fields
}
func EachQueryStructField(a interface{}, f func(jsonFieldName string, val reflect.Value, field reflect.StructField)) {
	EachSchemaField(a, func(jsonFieldName string, val reflect.Value, field reflect.StructField) {
		if field.Type.Kind() == reflect.Array || field.Type.Kind() == reflect.Slice {
			EachQueryStructField(val, f)
			return
		}
		f(jsonFieldName, val, field)
	})
}

// MapRows : sql.Rows を struct にマップ
// parentArray: *[]struct{}
func MapRows(parentArray interface{}, rows *sql.Rows) error {
	var columns []interface{}
	GetTypedColumns(parentArray, &columns, nil)
	slices, names, err := RowSlices(rows, columns)
	if err != nil {
		return err
	}
	for _, s := range slices {
		_, err := CreateOrUpdateMappedRow(parentArray, s, names, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateOrUpdateMappedRow : sql.Row を struct にマップ
// parentArray: *[]struct{}
func CreateOrUpdateMappedRow(parentArray interface{}, row []interface{}, columnNames []string, offset *int) (bool, error) {
	// 子の array も再帰検索
	// 原則上から

	// 数比較:  fields & columns
	// カラム名比較:  順番と一致しているか

	var parentArrayValue reflect.Value
	if v, ok := parentArray.(reflect.Value); ok {
		parentArrayValue = v
	} else {
		parentArrayValue = reflect.ValueOf(parentArray).Elem() // pointer indirect
	}
	childType := parentArrayValue.Type().Elem()

	var (
		columnFields  []string
		schemaFields  []int
		schemaColumns []interface{}
		keyFields     []int
		keyColumns    []interface{}
		childFields   []int
	)

	var columnNum = 0
	if offset == nil {
		offset = &columnNum
	}
	columnNum = *offset

	// field & column が一致しているか確かめる
	if columnNum == 0 {
		// TODO
		var allFields []string
		if len(allFields) != len(row) {
			// return false, fmt.Errorf("Field number mismatch: field (%d) != row (%d)", len(allFields), len(row))
		}
		if !cmp.Equal(allFields, row) {
			// return false, fmt.Errorf("Field name mismatch: field (%v) != row (%v)", allFields, row)
		}
	}

	var emptyChild = reflect.New(childType).Elem()
	EachSchemaField(emptyChild, func(jsonFieldName string, val reflect.Value, field reflect.StructField) {
		if columnNum >= len(row) {
			return
		}
		column := row[columnNum]

		tag := ParseGormTagSetting(field.Tag)
		if tag["PRIMARY_KEY"] != "" {
			keyFields = append(keyFields, columnNum)
			keyColumns = append(keyColumns, column)
		} else if field.Type.Kind() == reflect.Array || field.Type.Kind() == reflect.Slice {
			childFields = append(childFields, columnNum)
			return
		}

		var columnName string = tag["COLUMN"]
		if columnName == "" {
			columnName = ToSnakeCase(field.Name)
		}
		columnFields = append(columnFields, columnName)

		schemaFields = append(schemaFields, columnNum)
		schemaColumns = append(schemaColumns, column)

		columnNum++
	})

	// fmt.Println(row)
	// fmt.Println(columnNames)
	// fmt.Println(schemaFields)
	// fmt.Println(schemaColumns)
	// fmt.Println(keyFields)
	// fmt.Println(keyColumns)
	// fmt.Println(childFields)

	// 検索
	var newRow reflect.Value
	for i := 0; i < parentArrayValue.Len(); i++ {
		item := parentArrayValue.Index(i)

		found := true
		values := GetQueryStructValues(item)
		for ki, fieldNum := range keyFields {
			kfv := values[fieldNum]
			kf := kfv.Interface()
			ks := Convert(keyColumns[ki], kfv.Type())
			if !cmp.Equal(kf, ks) {
				found = false
				// fmt.Printf("not equal: %d : %#v != %#v, %v\n", i, kf, ks, kf == ks)
				break
			}
		}
		if found {
			// update
			newRow = item
			break
		}
	}

	// 変換
	var isNull = true
	var convertedColumns = make([]interface{}, len(schemaFields))
	values := GetQueryStructValues(emptyChild)
	for i, num := range schemaFields {
		field := values[num]
		val := reflect.ValueOf(schemaColumns[i])
		fmt.Println(field, num)

		if !val.IsValid() || val.IsZero() {
			// fmt.Println("nulalallalae", field.Type())
		} else {
			converted := Convert(val.Interface(), field.Type())
			if converted == nil {
				return false, fmt.Errorf("Field type mismatch: field %s -> %v (%v) != row -> %v (%v)", num, field, field.Type(), val, val.Type())
			}
			convertedColumns[i] = converted
			isNull = false
		}
	}

	if !isNull {
		if !newRow.IsValid() {
			// insert
			newRow = reflect.New(childType).Elem()
			newArray := reflect.Append(parentArrayValue, newRow)
			parentArrayValue.Set(newArray)
			newRow = parentArrayValue.Index(parentArrayValue.Len() - 1)
		}
		values := GetQueryStructValues(newRow)
		for i, num := range schemaFields {
			field := values[num]
			if convertedColumns[i] != nil {
				field.Set(reflect.ValueOf(convertedColumns[i]))
			}
		}
	}

	newRowValues := GetQueryStructValues(newRow)
	for _, cf := range childFields {
		f := newRowValues[cf]
		CreateOrUpdateMappedRow(f, row, columnNames, offset)
	}

	return true, nil
}

// RowSlices : rows -> [][]interface{} に変換 (typed)
func RowSlices(rows *sql.Rows, typedColumns []interface{}) ([][]interface{}, []string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	if len(columns) != len(typedColumns) {
		return nil, nil, fmt.Errorf("length mismatch: rows.Columns %d != typedColumns %d", len(columns), len(typedColumns))
	}
	rowList := [][]interface{}{}

	for rows.Next() {
		for i := range typedColumns {
			v := reflect.ValueOf(typedColumns[i]).Elem()
			v.Set(reflect.Zero(v.Type()))
		}
		rows.Scan(typedColumns...)
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = reflect.ValueOf(typedColumns[i]).Elem().Interface()
		}
		rowList = append(rowList, values)
	}
	rows.Close()

	return rowList, columns, nil
}

// GetTypedColumns : query struct -> []interface{} (typed)
func GetTypedColumns(parentArray interface{}, columns *[]interface{}, offset *int) error {
	var parentArrayValue reflect.Value
	if v, ok := parentArray.(reflect.Value); ok {
		parentArrayValue = v
	} else {
		parentArrayValue = reflect.ValueOf(parentArray).Elem() // pointer indirect
	}
	var childType = parentArrayValue.Type().Elem()
	var childFields []int

	var columnNum = 0
	if offset == nil {
		offset = &columnNum
	}
	columnNum = *offset

	var emptyChild = reflect.New(childType).Elem()
	EachSchemaField(emptyChild, func(jsonFieldName string, val reflect.Value, field reflect.StructField) {
		fmt.Println(field.Name)
		if field.Type.Kind() == reflect.Array || field.Type.Kind() == reflect.Slice {
			childFields = append(childFields, columnNum)
			return
		}
		newColumn := reflect.New(field.Type).Interface()
		(*columns) = append(*columns, newColumn)

		columnNum++
	})

	values := GetQueryStructValues(emptyChild)
	for _, cf := range childFields {
		f := values[cf]
		GetTypedColumns(f, columns, offset)
	}

	return nil
}

// AssignJSONFields
// *out = userStruct
func AssignJSONFields(out interface{}, userStruct interface{}) error {
	// fmt.Println("utils.AssignJSONFields:", reflect.TypeOf(out).Name(), reflect.TypeOf(userStruct).Name())

	var inMap Jsonmap = InterfaceToJsonmap(userStruct)
	if inMap == nil {
		inMap = Jsonmap{}
		StructToJsonmap(&inMap, userStruct)
	}
	var outMap Jsonmap = InterfaceToJsonmap(out)
	if outMap == nil {
		AssignMapToStruct(out, inMap)
	} else {
		for k, v := range inMap {
			outMap[k] = v
		}
	}

	return nil
}
func AssignJSONArrayFields(out interface{}, userArray interface{}) error {
	// fmt.Println("assignJSONArrayFields:", reflect.TypeOf(out).Name(), reflect.TypeOf(userArray).Name())
	uv := reflect.ValueOf(userArray)
	dv := reflect.Indirect(reflect.ValueOf(out))
	dElemType := dv.Type().Elem()

	for i := 0; i < uv.Len(); i++ {
		item := uv.Index(i)
		newItem := reflect.New(dElemType)
		err := AssignJSONFields(newItem.Interface(), item.Interface())
		if err != nil {
			return err
		}
		dv.Set(reflect.Append(dv, reflect.Indirect(newItem)))
	}
	return nil
}
