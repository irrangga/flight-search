package constant

func MapCabinClass(class string) string {
	switch class {
	case "F", "A":
		return "first"
	case "J", "C", "D", "Z":
		return "business"
	default:
		return "economy"
	}
}
