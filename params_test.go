package siesta

import (
	"net/url"
	"testing"
	"time"
)

func TestParamsSimple(t *testing.T) {
	v := url.Values{}
	p := Params{}
	v.Set("company", "VividCortex")
	v.Set("founded", "2012")
	v.Set("startup", "true")
	v.Set("duration", "10ms")
	v.Set("float", "12.89")
	v.Set("uint64", "1234")
	v.Set("int64", "-9876")
	v.Set("uint", "2345")
	v.Set("nonexistent", "8765")
	v.Set("valueless", "")
	v.Set("falseBool", "f")
	company := p.String("company", "", "the company name")
	founded := p.Int("founded", 0, "when it was founded")
	startup := p.Bool("startup", false, "whether it's a startup")
	duration := p.Duration("duration", 0, "how long it's been")
	floatVar := p.Float64("float", 0, "some float64")
	uint64Var := p.Uint64("uint64", 0, "some uint64")
	int64Var := p.Int64("int64", 0, "some int64")
	uintVar := p.Uint("uint", 0, "some uint")
	valueless := p.Bool("valueless", false, "some bool")
	falseBool := p.Bool("falseBool", true, "a bool with value false")
	err := p.Parse(v)
	if err != nil {
		t.Error(err)
	} else if *company != "VividCortex" {
		t.Errorf("expected VividCortex, got %s", *company)
	} else if *founded != 2012 {
		t.Errorf("expected 2012, got %d", *founded)
	} else if !*startup {
		t.Errorf("expected true, got %t", *startup)
	} else if *duration != 10*time.Millisecond {
		t.Errorf("expected 10ms, got %s", *duration)
	} else if *floatVar != 12.89 {
		t.Errorf("expected 12.89, got %f", *floatVar)
	} else if *uint64Var != 1234 {
		t.Errorf("expected 1234, got %d", *uint64Var)
	} else if *int64Var != -9876 {
		t.Errorf("expected -9876, got %d", *int64Var)
	} else if *uintVar != 2345 {
		t.Errorf("expected 2345, got %d", *uintVar)
	} else if *valueless != true {
		t.Errorf("expected true, got %t", *valueless)
	} else if *falseBool != false {
		t.Errorf("expected false, got %t", *falseBool)
	}

	usage := p.Usage()
	var expected map[string][3]string = map[string][3]string{
		"company":   [3]string{"company", "string", "the company name"},
		"founded":   [3]string{"founded", "int", "when it was founded"},
		"startup":   [3]string{"startup", "bool", "whether it's a startup"},
		"duration":  [3]string{"duration", "duration", "how long it's been"},
		"float":     [3]string{"float", "float64", "some float64"},
		"uint64":    [3]string{"uint64", "uint64", "some uint64"},
		"int64":     [3]string{"int64", "int64", "some int64"},
		"uint":      [3]string{"uint", "uint", "some uint"},
		"valueless": [3]string{"valueless", "bool", "some bool"},
		"falseBool": [3]string{"falseBool", "bool", "a bool with value false"},
	}
	compareUsageMaps(t, usage, expected)
}

func compareUsageMaps(t *testing.T, got, expected map[string][3]string) {
	seen := make(map[string]bool)
	for k, v := range got {
		seen[k] = true
		v2, ok := expected[k]
		if !ok {
			t.Errorf("%s doesn't exist in expected", k)
		} else if v2 != v {
			t.Errorf("%s: got '%s', expected '%s'", k, v, v2)
		}
	}
	for k, _ := range expected {
		if !seen[k] {
			_, ok := got[k]
			if !ok {
				t.Errorf("%s doesn't exist in got", k)
			}
		}
	}
}

func compareSlices(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, aVal := range a {
		if aVal != b[i] {
			return false
		}
	}
	return true
}

func TestParamsSlices(t *testing.T) {
	v := url.Values{}
	p := Params{}
	v.Add("company", "VividCortex")
	v.Add("company", "Inc,comma")
	company := p.SliceString("company", "", "the company name")
	v.Add("founded", "2012")
	v.Add("founded", "2012,2102")
	founded := p.SliceInt("founded", 0, "when it was founded")
	v.Add("startup", "true")
	v.Add("startup", "false,true")
	startup := p.SliceBool("startup", false, "whether it's a startup")
	v.Add("float", "12.89")
	v.Add("float", "1.25,12.625")
	floatVar := p.SliceFloat64("float", 0, "some float64")
	v.Add("uint64", "1234")
	v.Add("uint64", "1234,5678")
	uint64Var := p.SliceUint64("uint64", 0, "some uint64")
	v.Add("int64", "-9876")
	v.Add("int64", "-9876,8765")
	int64Var := p.SliceInt64("int64", 0, "some int64")
	v.Add("uint", "2345")
	v.Add("uint", "2345,3456")
	uintVar := p.SliceUint("uint", 0, "some uint")
	v.Add("duration", "10ms")
	v.Add("duration", "10s,12ms")
	duration := p.SliceDuration("duration", 0, "how long it's been")

	err := p.Parse(v)
	if err != nil {
		t.Error(err)
	}

	companies := []string{"VividCortex", "Inc", "comma"}
	for i, v := range *company {
		if v != companies[i] {
			t.Errorf("expected %s, got %s", companies[i], v)
		}
	}

	foundings := []int{2012, 2012, 2102}
	for i, v := range *founded {
		if v != foundings[i] {
			t.Errorf("expected %d, got %d", foundings[i], v)
		}
	}

	startups := []bool{true, false, true}
	for i, v := range *startup {
		if v != startups[i] {
			t.Errorf("expected %t, got %t", startups[i], v)
		}
	}

	floats := []float64{12.89, 1.25, 12.625}
	for i, v := range *floatVar {
		if v != floats[i] {
			t.Errorf("expected %f, got %f", floats[i], v)
		}
	}

	uint64s := []uint64{1234, 1234, 5678}
	for i, v := range *uint64Var {
		if v != uint64s[i] {
			t.Errorf("expected %d, got %d", uint64s[i], v)
		}
	}

	int64s := []int64{-9876, -9876, 8765}
	for i, v := range *int64Var {
		if v != int64s[i] {
			t.Errorf("expected %d, got %d", int64s[i], v)
		}
	}

	uints := []uint{2345, 2345, 3456}
	for i, v := range *uintVar {
		if v != uints[i] {
			t.Errorf("expected %d, got %d", uints[i], v)
		}
	}

	durations := []time.Duration{10 * time.Millisecond, 10 * time.Second, 12 * time.Millisecond}
	for i, v := range *duration {
		if v != durations[i] {
			t.Errorf("expected %s, got %s", durations[i], v)
		}
	}

	usage := p.Usage()
	var expected map[string][3]string = map[string][3]string{
		"company":  [3]string{"company", "[]string", "the company name"},
		"founded":  [3]string{"founded", "[]int", "when it was founded"},
		"startup":  [3]string{"startup", "[]bool", "whether it's a startup"},
		"duration": [3]string{"duration", "[]duration", "how long it's been"},
		"float":    [3]string{"float", "[]float64", "some float64"},
		"uint64":   [3]string{"uint64", "[]uint64", "some uint64"},
		"int64":    [3]string{"int64", "[]int64", "some int64"},
		"uint":     [3]string{"uint", "[]uint", "some uint"},
	}
	compareUsageMaps(t, usage, expected)

}
