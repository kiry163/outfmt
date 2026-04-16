package outfmt

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type tableRows struct {
	headers []string
	rows    [][]string
}

type tableColumn struct {
	key    string
	header string
}

type flattenState struct {
	columns []tableColumn
	seen    map[string]struct{}
	row     map[string]string
}

var (
	stringerType      = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	yamlMarshalerType = reflect.TypeOf((*yaml.Marshaler)(nil)).Elem()
)

func extractRows(data any) (tableRows, error) {
	if data == nil {
		return tableRows{}, nil
	}

	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return tableRows{}, nil
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Struct:
		return extractSingleRow(value)
	case reflect.Map:
		return extractSingleRow(value)
	case reflect.Slice, reflect.Array:
		return extractSliceRows(value)
	default:
		return tableRows{}, fmt.Errorf("outfmt: table format does not support %s", value.Kind())
	}
}

func extractSingleRow(value reflect.Value) (tableRows, error) {
	columns, cells, err := flattenRow(value)
	if err != nil {
		return tableRows{}, err
	}
	return materializeRows(columns, []map[string]string{cells}), nil
}

func extractSliceRows(value reflect.Value) (tableRows, error) {
	if value.Len() == 0 {
		return tableRows{}, nil
	}

	elemType := stripPointerType(value.Type().Elem())
	switch elemType.Kind() {
	case reflect.Struct, reflect.Map:
	default:
		return tableRows{}, fmt.Errorf("outfmt: table slice elements must be struct or map, got %s", elemType.Kind())
	}

	var columns []tableColumn
	rows := make([]map[string]string, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		item, err := normalizeSliceItem(value.Index(i), elemType)
		if err != nil {
			return tableRows{}, err
		}

		itemColumns, cells, err := flattenRow(item)
		if err != nil {
			return tableRows{}, err
		}

		columns = mergeColumns(columns, itemColumns)
		rows = append(rows, cells)
	}

	return materializeRows(columns, rows), nil
}

func flattenRow(value reflect.Value) ([]tableColumn, map[string]string, error) {
	state := &flattenState{
		seen: make(map[string]struct{}),
		row:  make(map[string]string),
	}

	if err := flattenValue(value, nil, nil, state); err != nil {
		return nil, nil, err
	}

	return state.columns, state.row, nil
}

func flattenValue(value reflect.Value, path []string, headers []string, state *flattenState) error {
	if !value.IsValid() {
		state.setCell(path, headers, "")
		return nil
	}

	for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.Kind() == reflect.Interface {
			if value.IsNil() {
				state.setCell(path, headers, "")
				return nil
			}
			value = value.Elem()
			continue
		}

		if value.IsNil() {
			elemType := stripPointerType(value.Type().Elem())
			if elemType.Kind() == reflect.Struct && shouldFlattenStructType(elemType) {
				return flattenEmptyStruct(elemType, path, headers, state)
			}
			state.setCell(path, headers, "")
			return nil
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Struct:
		if shouldFlattenStructType(value.Type()) {
			return flattenStruct(value, path, headers, state)
		}
	case reflect.Map:
		if value.Type().Key().Kind() != reflect.String {
			return fmt.Errorf("outfmt: table map keys must be strings")
		}
		return flattenMap(value, path, headers, state)
	}

	state.setCell(path, headers, valueToCell(value))
	return nil
}

func flattenEmptyStruct(typ reflect.Type, path []string, headers []string, state *flattenState) error {
	typ = stripPointerType(typ)
	added := false

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}

		header, skip := tableHeader(field)
		if skip {
			continue
		}

		added = true
		childPath := appendPath(path, field.Name)
		childHeader := appendPath(headers, header)
		fieldType := stripPointerType(field.Type)

		if fieldType.Kind() == reflect.Struct && shouldFlattenStructType(fieldType) {
			if err := flattenEmptyStruct(fieldType, childPath, childHeader, state); err != nil {
				return err
			}
			continue
		}

		state.setCell(childPath, childHeader, "")
	}

	if !added {
		state.setCell(path, headers, "")
	}

	return nil
}

func flattenStruct(value reflect.Value, path []string, headers []string, state *flattenState) error {
	typ := value.Type()
	added := false
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}

		header, skip := tableHeader(field)
		if skip {
			continue
		}

		added = true
		if err := flattenValue(
			value.Field(i),
			appendPath(path, field.Name),
			appendPath(headers, header),
			state,
		); err != nil {
			return err
		}
	}

	if !added {
		state.setCell(path, headers, valueToCell(value))
	}

	return nil
}

func flattenMap(value reflect.Value, path []string, headers []string, state *flattenState) error {
	keys := value.MapKeys()
	if len(keys) == 0 {
		return nil
	}

	names := make([]string, 0, len(keys))
	for _, key := range keys {
		if key.Kind() != reflect.String {
			return fmt.Errorf("outfmt: table map keys must be strings")
		}
		names = append(names, key.String())
	}
	sort.Strings(names)

	for _, name := range names {
		if err := flattenValue(
			value.MapIndex(reflect.ValueOf(name)),
			appendPath(path, name),
			appendPath(headers, name),
			state,
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *flattenState) setCell(path []string, headers []string, cell string) {
	columnKey, header := normalizeColumn(path, headers)
	if columnKey == "" {
		return
	}

	if _, ok := s.seen[columnKey]; !ok {
		s.columns = append(s.columns, tableColumn{
			key:    columnKey,
			header: header,
		})
		s.seen[columnKey] = struct{}{}
	}
	s.row[columnKey] = cell
}

func normalizeColumn(path []string, headers []string) (string, string) {
	if len(path) == 0 {
		return "value", "Value"
	}
	return strings.Join(path, "."), strings.Join(headers, ".")
}

func tableHeader(field reflect.StructField) (string, bool) {
	tag := strings.TrimSpace(field.Tag.Get("outfmt"))
	if tag == "-" {
		return "", true
	}
	if tag == "" {
		return field.Name, false
	}
	return tag, false
}

func materializeRows(columns []tableColumn, flatRows []map[string]string) tableRows {
	headers := make([]string, 0, len(columns))
	for _, col := range columns {
		headers = append(headers, col.header)
	}

	rows := make([][]string, 0, len(flatRows))
	for _, flatRow := range flatRows {
		row := make([]string, 0, len(columns))
		for _, col := range columns {
			row = append(row, flatRow[col.key])
		}
		rows = append(rows, row)
	}

	return tableRows{
		headers: headers,
		rows:    rows,
	}
}

func mergeColumns(base []tableColumn, added []tableColumn) []tableColumn {
	if len(added) == 0 {
		return base
	}

	seen := make(map[string]struct{}, len(base))
	for _, col := range base {
		seen[col.key] = struct{}{}
	}
	for _, col := range added {
		if _, ok := seen[col.key]; ok {
			continue
		}
		base = append(base, col)
		seen[col.key] = struct{}{}
	}
	return base
}

func normalizeSliceItem(value reflect.Value, elemType reflect.Type) (reflect.Value, error) {
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Zero(elemType), nil
		}
		value = value.Elem()
	}

	if value.Kind() != elemType.Kind() {
		return reflect.Value{}, fmt.Errorf("outfmt: mixed slice element kinds are not supported")
	}

	return value, nil
}

func stripPointerType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	return typ
}

func appendPath(path []string, item string) []string {
	out := make([]string, 0, len(path)+1)
	out = append(out, path...)
	out = append(out, item)
	return out
}

func shouldFlattenStructType(typ reflect.Type) bool {
	typ = stripPointerType(typ)
	if typ.Kind() != reflect.Struct {
		return false
	}
	if implementsLeafInterface(typ) {
		return false
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}
		if field.Tag.Get("outfmt") == "-" {
			continue
		}
		return true
	}
	return false
}

func implementsLeafInterface(typ reflect.Type) bool {
	return typeImplements(typ, stringerType) ||
		typeImplements(typ, textMarshalerType) ||
		typeImplements(typ, jsonMarshalerType) ||
		typeImplements(typ, yamlMarshalerType)
}

func typeImplements(typ reflect.Type, iface reflect.Type) bool {
	if typ.Implements(iface) {
		return true
	}
	if typ.Kind() != reflect.Pointer && reflect.PointerTo(typ).Implements(iface) {
		return true
	}
	return false
}

func valueToCell(value reflect.Value) string {
	if !value.IsValid() {
		return ""
	}

	for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.Kind() == reflect.Interface {
			if value.IsNil() {
				return ""
			}
			value = value.Elem()
			continue
		}

		if value.IsNil() {
			return ""
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Bool:
		if value.Bool() {
			return "true"
		}
		return "false"
	case reflect.Map, reflect.Slice:
		if value.IsNil() {
			return ""
		}
		return fmt.Sprint(value.Interface())
	default:
		return fmt.Sprint(value.Interface())
	}
}
