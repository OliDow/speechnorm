package speechnorm

// intToBase10 is a shared test helper that formats an int64 as its base-10
// string representation, used to generate stable subtest names in
// locale-specific test tables.
func intToBase10(n int64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

// englishName is used for subtest names in English tests, including
// negative numbers which become "neg_<abs>".
func englishName(n int64) string {
	if n < 0 {
		return "neg_" + englishName(-n)
	}
	return intToBase10(n)
}
