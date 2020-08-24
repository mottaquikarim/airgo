package airgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	airtable "github.com/mottaquikarim/go-airtable"
)

type RebuildSiteOpts struct {
	Datadir       string
	Template      string
	Layout        string
	Output        string
	TemplateNames []string
}

func RebuildSite(cfg RebuildSiteOpts) error {
	// read all files in director
	var files []os.FileInfo
	var err error

	if len(cfg.Datadir) == 0 {
		return fmt.Errorf("Datadir path is required")
	} else if len(cfg.Template) == 0 {
		return fmt.Errorf("Template path is required")
	} else if len(cfg.Output) == 0 {
		return fmt.Errorf("Output path is required")
	} else if len(cfg.TemplateNames) == 0 {
		return fmt.Errorf("At least one templateName required")
	}

	if len(cfg.Layout) == 0 {
		cfg.Layout = "layout.html.tmpl"
	}

	if files, err = ioutil.ReadDir(cfg.Datadir); err != nil {
		return fmt.Errorf("Failed to read files from dir %v %v", cfg.Datadir, err)
	}

	layoutTmpl := fmt.Sprintf("%s/%s", cfg.Template, cfg.Layout)
	for _, f := range files {
		var content []byte
		var filename string
		var fileRoot string

		filename = fmt.Sprintf("%s/%s", cfg.Datadir, f.Name())
		if content, err = ioutil.ReadFile(filename); err != nil {
			return fmt.Errorf("Failed to read file: %s %v", filename, err)
		}

		fileRoot = strings.Split(f.Name(), ".")[0]

		// marshal content into airtable.Record
		var data = []airtable.Record{}
		if err = json.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("Failed to unmarshal json %v", err)
		}

		// write file through template
		// more info here:
		// https://www.htmlgoodies.com/beyond/reference/nesting-templates-with-go-web-programming.html
		templateName := fmt.Sprintf("%s/%s.html.tmpl", cfg.Template, fileRoot)
		t, err := template.ParseFiles(layoutTmpl, templateName)
		if err != nil {
			return fmt.Errorf("Failed to parse file: %v", err)
		}

		// create or truncate index.html
		f, err := os.Create(fmt.Sprintf("%s/%s.html", cfg.Output, fileRoot))
		defer f.Close()
		if err != nil {
			return fmt.Errorf("Failed to create file: %v", err)
		}

		// populate template with data
		for _, tmplName := range cfg.TemplateNames {
			if err = t.ExecuteTemplate(f, tmplName, data); err != nil {
				return fmt.Errorf("Failed to execute: %v", err)
			}
		}
	}

	return nil
}
