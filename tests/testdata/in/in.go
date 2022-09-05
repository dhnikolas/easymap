package in

type Source struct {
	Name      string `json:"name"`
	Age       int
	nonpublic string
	Docs      *Doci
	Users     []uint16
	Disks     []*MyDisk
	Meta      map[string]string
}

type Doci struct {
	Serial string
	Number string
	Rights *MyRights
}

type MyRights struct {
	Some string
	Laws []*MyLaw
}

type MyLaw struct {
	State  string
	Params []*MyParam
}

type MyParam struct {
	Count int
}

type MyDisk struct {
	Size        int
	Description string
	Zone        []*AvailabilityZone
}

type AvailabilityZone struct {
	Number int
	City   string
	Config []*MyOptions
	MTU    string
}

type MyOptions struct {
	TTL     string
	Timeout int
	Simple  []*MySimple
}

type MySimple struct {
	Main       string
	Additional MyAdd
}

type MyAdd struct {
	Column string
	Rows   []*MyRow
}

type MyRow struct {
	Number int
}

var TestData = map[string]string{
	"Name":                            "test1",
	"Docs_Serial":                     "serial",
	"Docs_Number":                     "555",
	"Docs_Rights_Some":                "some",
	"Docs_Rights_Laws1_State":         "law1",
	"Docs_Rights_Laws2_State":         "law2",
	"Disks_Description":               "description",
	"Disks_Zone1_City":                "city",
	"Disks_Zone1_Mtu":                 "mtu",
	"Disks_Zone2_City":                "city2",
	"Disks_Zone2_Mtu":                 "mtu2",
	"Disks_Zone2_Config_TTL":          "ttl",
	"Disks_Zone2_Config_Simple1_Main": "main",
	"Disks_Zone2_Config_Simple1_Additional_Column": "col",
}

func GenInTestStruct() *Source {
	src := &Source{
		Name:      TestData["Name"],
		Age:       2,
		nonpublic: "",
		Docs: &Doci{
			Serial: TestData["Docs_Serial"],
			Number: TestData["Docs_Number"],
			Rights: &MyRights{
				Some: TestData["Docs_Rights_Some"],
				Laws: []*MyLaw{
					&MyLaw{
						State:  TestData["Docs_Rights_Laws1_State"],
						Params: []*MyParam{&MyParam{Count: 2}},
					}, &MyLaw{
						State:  TestData["Docs_Rights_Laws2_State"],
						Params: []*MyParam{&MyParam{Count: 2}},
					},
				},
			},
		},
		Users: []uint16{2},
		Disks: []*MyDisk{
			&MyDisk{
				Size:        2,
				Description: TestData["Disks_Description"],
				Zone: []*AvailabilityZone{
					&AvailabilityZone{
						Number: 2,
						City:   TestData["Disks_Zone1_City"],
						Config: nil,
						MTU:    TestData["Disks_Zone1_Mtu"],
					},
					&AvailabilityZone{
						Number: 2,
						City:   TestData["Disks_Zone2_City"],
						Config: []*MyOptions{
							{
								TTL:     TestData["Disks_Zone2_Config_TTL"],
								Timeout: 2,
								Simple: []*MySimple{
									{
										Main: TestData["Disks_Zone2_Config_Simple1_Main"],
										Additional: MyAdd{
											Column: TestData["Disks_Zone2_Config_Simple1_Additional_Column"],
											Rows:   []*MyRow{{Number: 2}},
										},
									},
								},
							},
						},
						MTU: TestData["Disks_Zone2_Mtu"],
					},
				},
			},
		},
		Meta: map[string]string{"2": "2"},
	}
	return src
}
