package typeconv

func DbToGo(db string) string {
	switch db {
	case "string", "text":
		return "string"
	case "int", "serial":
		return "int"
	case "bool":
		return "bool"
	case "json":
		return "json.RawMessage"
	case "uuid":
		return "uuid.UUID"
	case "timestamp(0)":
		return "time.Time"
	default:
		return ""
	}
}

func BoolToString(v bool) string {
	if v {
		return "1"
	}
	return ""
}
