package cast


func CastToString(s string) *string{
	if s == ""{
		return nil
	}
	return &s
}