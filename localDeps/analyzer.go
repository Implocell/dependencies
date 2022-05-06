package localdeps

type EmptyExport struct {
	FileName      string   `json:"filename"`
	UnusedExports []string `json:"unusedExports"`
}

func EmptyExports(fileDeps []FileDeps) []EmptyExport {
	var emptyExports []EmptyExport

	for _, f := range fileDeps {
		var tmpEmptyExports []string
		for _, e := range f.Exports {
			if len(e.UsedBy) == 0 {
				tmpEmptyExports = append(tmpEmptyExports, e.Export)
			}
		}
		if len(tmpEmptyExports) > 0 {
			emptyExports = append(emptyExports,
				EmptyExport{
					FileName:      f.Filename,
					UnusedExports: tmpEmptyExports,
				},
			)
		}
	}

	return emptyExports
}
