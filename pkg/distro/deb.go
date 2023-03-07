package distro

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	spdx "github.com/spdx/tools-golang/spdx/v2/v2_3"
	"pault.ag/go/debian/deb"
	k8spdx "sigs.k8s.io/bom/pkg/spdx"
)

func ParseFile(path string) error {
	debfile, _, err := deb.LoadFile(path)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		// FIXME:
	}

	defer debfile.Close()

	kdoc := k8spdx.NewDocument()
	pkg := &k8spdx.Package{
		Entity: k8spdx.Entity{
			Name: debfile.Control.Package,
		},
		Version: debfile.Control.Version.String(),
		Originator: struct {
			Person       string
			Organization string
		}{
			Person: debfile.Control.Maintainer,
		},
	}
	kdoc.AddPackage(pkg)

	//	fmt.Printf("deb: %+v\n", debfile)

	/*
		for k, v := range debfile.ArContent {
			fmt.Printf("KEY: %+v VAL: %+v\n", k, v)
			fmt.Printf("LICENSE found\n")
			firstContent, _ := ioutil.ReadAll(v.Data)
			fmt.Printf("firstContent:%+v\n", string(firstContent))
		}
	*/

	var files []*spdx.File = []*spdx.File{}
	for {
		hdr, err := debfile.Data.Next()
		if hdr == nil {
			break
		}

		buf := make([]byte, hdr.Size)
		debfile.Data.Read(buf)

		files = append(files, &spdx.File{FileName: hdr.Name[1:]})
		pkg.AddFile(&k8spdx.File{Entity: k8spdx.Entity{Name: hdr.Name[1:]}})

		if strings.HasPrefix(hdr.Name, "./usr/share/doc/") && strings.HasSuffix(hdr.Name, "copyright") {
			fmt.Printf("LICENSE found\n")
			//fmt.Printf("%s", string(buf))
			pkg.LicenseComments = string(buf)
		}

		if err != nil {
			fmt.Printf("err: %+v\n", err)
		}

		//		fmt.Printf("hdr: +%v\n", hdr)
	}

	doc := spdx.Document{
		SPDXVersion: "SPDX-2.3",
		DataLicense: "CC0-1.0",
		Packages: []*spdx.Package{
			&spdx.Package{PackageName: debfile.Control.Package, Files: files}},
	}

	fmt.Printf("%+v\n", doc)

	bytes, err := json.Marshal(doc)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}

	ioutil.WriteFile(path+".spdx", bytes, 0644)
	kdoc.Write(path + ".k8s.spdx")

	return nil
}
