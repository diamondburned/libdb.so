package resume

import (
	"encoding/json"
	"fmt"
	"go/doc"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"
	"libdb.so/vm"
	"libdb.so/vm/internal/ansi"
	"libdb.so/vm/internal/cliprog"
	"libdb.so/vm/programs"
)

func init() {
	programs.Register(cliprog.Wrap(app))
}

// Document is the main struct for the resume.
type Document struct {
	Awards           []Award     `json:"awards"`
	Basics           Basics      `json:"basics"`
	Education        []Education `json:"education"`
	Projects         []Project   `json:"projects"`
	SelectedTemplate int64       `json:"selectedTemplate"`
	Skills           []Skill     `json:"skills"`
	Work             []Work      `json:"work"`
	Sections         []string    `json:"sections"`
}

// Education is the struct for the education section.
type Education struct {
	Area        string `json:"area"`
	EndDate     string `json:"endDate"`
	Institution string `json:"institution"`
	Location    string `json:"location"`
	StartDate   string `json:"startDate"`
	StudyType   string `json:"studyType"`
}

// Award is the struct for the awards section.
type Award struct {
	Awarder string `json:"awarder"`
	Date    string `json:"date"`
	Summary string `json:"summary"`
	Title   string `json:"title"`
}

// Work is the struct for the work section.
type Work struct {
	Company    string   `json:"company"`
	EndDate    string   `json:"endDate"`
	Highlights []string `json:"highlights"`
	Location   string   `json:"location"`
	Position   string   `json:"position"`
	StartDate  string   `json:"startDate"`
	Website    string   `json:"website,omitempty"`
}

// Project is the struct for the projects section.
type Project struct {
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Name        string   `json:"name"`
	URL         string   `json:"url"`
}

// Basics is the struct for the basics section.
type Basics struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Website string `json:"website"`
}

// Skill is the struct for the skills section.
type Skill struct {
	Keywords []string `json:"keywords"`
	Name     string   `json:"name"`
	Level    string   `json:"level,omitempty"`
}

var app = cli.App{
	Name:  "resume",
	Usage: "print resume",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "print resume in json format",
		},
		&cli.IntFlag{
			Name:    "width",
			Aliases: []string{"w"},
			Usage:   "set the width of the printing viewport",
		},
	},
	Action: action,
}

func action(c *cli.Context) error {
	env := vm.EnvironmentFromContext(c.Context)

	width := env.Terminal.Query().Width
	if c.IsSet("width") {
		width = c.Int("width")
	}

	const maxWidth = 80
	const minWidth = 15
	if width > maxWidth {
		width = maxWidth
	}
	if width < minWidth {
		width = minWidth
	}

	f, err := env.Open("/resume.json")
	if err != nil {
		return errors.Wrap(err, "failed to open resume.json")
	}
	defer f.Close()

	if c.Bool("json") {
		_, err := io.Copy(env.Terminal.Stdout, f)
		return errors.Wrap(err, "failed to copy resume.json to stdout")
	}

	var doc Document
	if err := json.NewDecoder(f).Decode(&doc); err != nil {
		return errors.Wrap(err, "failed to decode resume.json")
	}

	const name = "diamondburned"
	ansi.PrintAligned(env, width,
		"",
		spcolor(
			"━━━ "+name+" ━━━",
			color.Bold, color.FgHiMagenta),
		fmt.Sprintf("%s  %s",
			ansi.Link(doc.Basics.Email, "mailto://"+doc.Basics.Email),
			ansi.Link(doc.Basics.Website, "https://"+doc.Basics.Website)))
	env.Println()

	sectionPrinters := map[string]func(){
		"education": func() {
			if len(doc.Education) == 0 {
				return
			}

			pheading(env, width, "Education")
			for _, education := range doc.Education {
				ansi.PrintAligned(env, width-2,
					"  "+spcolor(education.Institution, color.Bold), "",
					"  "+spcolor(education.Area, color.Italic))
				ansi.PrintAligned(env, width-2,
					"  "+education.Area, "", fmt.Sprintf("%s - %s", education.StartDate, education.EndDate))
				env.Println()
			}
		},
		"work": func() {
			if len(doc.Work) == 0 {
				return
			}

			pheading(env, width, "Work Experiences")
			for _, work := range doc.Work {
				ansi.PrintAligned(env, width-2,
					fmt.Sprintf("  %s /%s/",
						spcolor(work.Company, color.Bold),
						spcolor(work.Position, color.Italic)), "",
					fmt.Sprintf("  %s - %s", work.StartDate, work.EndDate))
				for _, highlight := range work.Highlights {
					env.Println(wrap("- "+highlight, width-2, 2))
				}
				env.Println()
			}
		},
		"skills": func() {
			if len(doc.Skills) == 0 {
				return
			}

			pheading(env, width, "Skills")
			for _, skill := range doc.Skills {
				env.Println("  "+spcolor(skill.Name, color.Bold), strings.Join(skill.Keywords, ", "))
			}
			env.Println()
		},
		"projects": func() {
			if len(doc.Projects) == 0 {
				return
			}

			pheading(env, width, "Projects")
			for _, project := range doc.Projects {
				ansi.PrintAligned(env, width-2,
					fmt.Sprintf("  %s  %s",
						spcolor(project.Name, color.Bold),
						strings.Join(project.Keywords, ", ")), "",
					ansi.Link(project.URL, project.URL))
				env.Println(wrap(project.Description, width-2, 2))
				env.Println()
			}
		},
	}

	filtered := doc.Sections[:0]
	for _, section := range doc.Sections {
		if _, ok := sectionPrinters[section]; ok {
			filtered = append(filtered, section)
		}
	}

	doc.Sections = filtered
	for i, section := range doc.Sections {
		printer := sectionPrinters[section]
		printer()

		if i != len(doc.Sections)-1 {
			env.Print("Press Enter to continue...")
			env.Terminal.Stdin.Read(make([]byte, 1))
			env.Print(ansi.ClearLine, ansi.MoveCursorToStart)
		}
	}

	return nil
}

func spcolor(s string, cs ...color.Attribute) string {
	return color.New(cs...).Sprint(s)
}

func pcolorln(env vm.Environment, s string, cs ...color.Attribute) {
	fmt.Fprintln(env.Terminal.Stdout, color.New(cs...).Sprint(s))
}

func pheading(env vm.Environment, width int, heading string) {
	const fgcolor = color.FgHiGreen

	w := ansi.StringWidth(heading)
	ansi.PrintAligned(env, width,
		fmt.Sprint(
			spcolor("╭─", fgcolor, color.Bold),
			spcolor(heading, fgcolor, color.Bold),
			spcolor(strings.Repeat("─", padwidth(w)), fgcolor, color.Bold),
		), "",
		spcolor(" ─╮", fgcolor, color.Bold),
	)
}

func wrap(line string, width int, pad int) string {
	var out strings.Builder
	doc.ToText(&out, line, strings.Repeat(" ", pad), "", width-pad)
	return strings.TrimRight(out.String(), "\n")
}

func padwidth(strwidth int) int {
	return max(30-2-strwidth, 0)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
