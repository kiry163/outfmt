package outfmt

import "testing"

type sampleUser struct {
	ID     int    `json:"id" yaml:"id" outfmt:"ID"`
	Name   string `json:"name" yaml:"name"`
	Email  string `json:"email,omitempty" yaml:"email,omitempty"`
	Secret string `outfmt:"-"`
}

func TestMarshalJSON(t *testing.T) {
	raw, err := Marshal(sampleUser{
		ID:     1,
		Name:   "alice",
		Email:  "alice@example.com",
		Secret: "ignored",
	}, JSON)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	want := "{\n  \"id\": 1,\n  \"name\": \"alice\",\n  \"email\": \"alice@example.com\",\n  \"Secret\": \"ignored\"\n}\n"
	if string(raw) != want {
		t.Fatalf("unexpected json output:\nwant:\n%s\ngot:\n%s", want, string(raw))
	}
}

func TestMarshalYAML(t *testing.T) {
	raw, err := Marshal(sampleUser{
		ID:    1,
		Name:  "alice",
		Email: "alice@example.com",
	}, YAML)
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	want := "id: 1\nname: alice\nemail: alice@example.com\nsecret: \"\"\n"
	if string(raw) != want {
		t.Fatalf("unexpected yaml output:\nwant:\n%s\ngot:\n%s", want, string(raw))
	}
}

func TestMarshalTableFromStructSlice(t *testing.T) {
	got := mustWriteTable([]sampleUser{
		{ID: 1, Name: "alice", Email: "alice@example.com", Secret: "ignored"},
		{ID: 2, Name: "bob"},
	})

	want := "ID  Name   Email            \n--  -----  -----------------\n1   alice  alice@example.com\n2   bob    -                \n"
	if got != want {
		t.Fatalf("unexpected table output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestMarshalTableFromMapSlice(t *testing.T) {
	got := mustWriteTable([]map[string]any{
		{"name": "alice", "score": 10},
		{"name": "bob"},
	})

	want := "name   score\n-----  -----\nalice  10   \nbob    -    \n"
	if got != want {
		t.Fatalf("unexpected table output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

type sampleProfile struct {
	City   string `outfmt:"City"`
	Active bool   `outfmt:"Active"`
}

type sampleNestedUser struct {
	ID      int            `outfmt:"ID"`
	Profile *sampleProfile `outfmt:"Profile"`
	Meta    map[string]any `outfmt:"Meta"`
	Tags    []string       `outfmt:"Tags"`
}

func TestMarshalTableFromNestedStructSlice(t *testing.T) {
	got := mustWriteTable([]sampleNestedUser{
		{
			ID:      1,
			Profile: &sampleProfile{City: "shanghai", Active: true},
			Meta:    map[string]any{"region": "cn", "zone": "east"},
			Tags:    []string{"dev", "ops"},
		},
		{
			ID:      2,
			Profile: nil,
			Meta:    map[string]any{"region": "us"},
		},
	})

	want := "ID  Profile.City  Profile.Active  Meta.region  Meta.zone  Tags     \n--  ------------  --------------  -----------  ---------  ---------\n1   shanghai      true            cn           east       [dev ops]\n2   -             -               us           -          -        \n"
	if got != want {
		t.Fatalf("unexpected nested struct table output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestMarshalTableFromNestedMapSlice(t *testing.T) {
	got := mustWriteTable([]map[string]any{
		{
			"name": "alice",
			"profile": map[string]any{
				"age":  18,
				"city": "shanghai",
			},
		},
		{
			"name": "bob",
			"profile": map[string]any{
				"city": "beijing",
			},
		},
	})

	want := "name   profile.age  profile.city\n-----  -----------  ------------\nalice  18           shanghai    \nbob    -            beijing     \n"
	if got != want {
		t.Fatalf("unexpected nested map table output:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestMarshalTableUnsupportedType(t *testing.T) {
	_, err := Marshal("plain-text", Table)
	if err == nil {
		t.Fatal("expected error for unsupported table type")
	}
}
