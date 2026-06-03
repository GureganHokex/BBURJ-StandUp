package admin

type Model struct {
	Name     string
	Slug     string
	APIPath  string
	ListPath string
}

func Models() []Model {
	return []Model{
		{Name: "Site Settings", Slug: "settings", APIPath: "/api/settings", ListPath: "/admin/settings"},
		{Name: "Events", Slug: "events", APIPath: "/api/events", ListPath: "/admin/events"},
		{Name: "Videos", Slug: "videos", APIPath: "/api/videos", ListPath: "/admin/videos"},
		{Name: "Photos", Slug: "photos", APIPath: "/api/photos", ListPath: "/admin/photos"},
		{Name: "Merch", Slug: "merch", APIPath: "/api/merch", ListPath: "/admin/merch"},
	}
}

func FindBySlug(slug string) (Model, bool) {
	for _, m := range Models() {
		if m.Slug == slug {
			return m, true
		}
	}
	return Model{}, false
}
