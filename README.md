# Map your golang struct really easy

The library is designed for fast mapping of structs in the Go language.
You can transfer data from one structure to another using a generated method. Assigning fields with the same name

```
$ easymap gen /path/to/go/file.go:StructName /path/to/another/file.go:AnotherStruct

func MapToAnotherStruct(in *pb.StructName) *AnotherStruct {
	out := &AnotherStruct{
		ProjectId:        in.ProjectId,
		Name:             in.Name,
		FlavorName:       in.FlavorName,
		ImageName:        in.ImageName,
		NetworkName:      in.NetworkName,
		SecurityGroupsId: in.SecurityGroupsId,
		Tags:             in.Tags,
	}
	if in.Keys != nil {
		out.Keys = &Keys{
			KeyName: in.Keys.KeyName,
		}
		if in.Keys.New != nil {
			out.Keys.New = &Keypair{
				Name:      in.Keys.New.Name,
				Type:      in.Keys.New.Type,
				PublicKey: in.Keys.New.PublicKey,
				UserId:    in.Keys.New.UserId,
			}
		}
	}
	if in.Disks != nil {
		out.Disks = &VMDisks{}
		var newDisksIdSlice []*DisksId
		for _, DisksIdItem := range in.Disks.DisksId {
			newDisksId := &DisksId{
				Id:   DisksIdItem.Id,
				Boot: DisksIdItem.Boot,
			}
			newDisksIdSlice = append(newDisksIdSlice, newDisksId)
		}
		out.Disks.DisksId = newDisksIdSlice
		var newNewSlice []*NewDisk
		for _, NewItem := range in.Disks.New {
			newNew := &NewDisk{
				DisplayName:        NewItem.DisplayName,
				DisplayDescription: NewItem.DisplayDescription,
				Size:               NewItem.Size,
				TypeId:             NewItem.TypeId,
				Boot:               NewItem.Boot,
			}
			newNewSlice = append(newNewSlice, newNew)
		}
		out.Disks.New = newNewSlice
	}
	if in.FloatingIp != nil {
		out.FloatingIp = &CreateVM_Request_FloatingIp{
			AllocateNew: in.FloatingIp.AllocateNew,
			Description: in.FloatingIp.Description,
		}
	}
	return out
}

```

It is also possible to generate a dto structure based on another struct and a mapping method between them
```
$ easymap copygen /example/protocols/pb/integration_api_sg.pb.go:SG NewSG

type NewSG struct {
	Id		string		
	Name		string		
	Description	string		
	Tags		[]string	
	Rules		[]*SecGroupRule	
}
type SecGroupRule struct {
	Id		string	
	Direction	string	
	Description	string	
	EtherType	string	
	SecGroupId	string	
	PortRangeMin	int32	
	PortRangeMax	int32	
	Protocol	string	
	RemoteGroupId	string	
	RemoteIpPrefix	string	
	TenantId	string	
	ProjectId	string	
}


func MapToNewSG(in *pb.SG) *NewSG {
	out := &NewSG{
		Id:          in.Id,
		Name:        in.Name,
		Description: in.Description,
		Tags:        in.Tags,
	}
	var newRulesSlice []*SecGroupRule
	for _, RulesItem := range in.Rules {
		newRules := &SecGroupRule{
			Id:             RulesItem.Id,
			Direction:      RulesItem.Direction,
			Description:    RulesItem.Description,
			EtherType:      RulesItem.EtherType,
			SecGroupId:     RulesItem.SecGroupId,
			PortRangeMin:   RulesItem.PortRangeMin,
			PortRangeMax:   RulesItem.PortRangeMax,
			Protocol:       RulesItem.Protocol,
			RemoteGroupId:  RulesItem.RemoteGroupId,
			RemoteIpPrefix: RulesItem.RemoteIpPrefix,
			TenantId:       RulesItem.TenantId,
			ProjectId:      RulesItem.ProjectId,
		}
		newRulesSlice = append(newRulesSlice, newRules)
	}
	out.Rules = newRulesSlice
	return out
}
```


### Install on Mac OS with brew
```
 brew install dhnikolas/tools/easymap
```

### Install on Linux
```
$ curl -s https://api.github.com/repos/dhnikolas/easymap/releases \
| grep "browser_download_url.*linux.tar.gz" | head -n1 \
| cut -d : -f 2,3 | tr -d \" | wget -qi - -O -| sudo tar -xz -C /usr/local/bin
 
$ sudo chmod +x  /usr/local/bin/easymap
```

Happy mapping :smile: