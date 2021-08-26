package upwork

import "testing"

func TestParseSummary(t *testing.T) {
	type testValue struct {
		Src string
		Sum Summary
	}

	tests := []testValue{
		{
			Src: "Description<br /><br /><br /><b>Posted On</b>: August 24, 2021 07:28 UTC<br /><b>Category</b>: Front-End Development<br /><b>Skills</b>:React Bootstrap,     Chart.js\n<br /><b>Country</b>: United States\n<br /><a href=\"https://www.upwork.com/jobs/?source=rss\">click to apply</a>\n",
			Sum: Summary{
				Country: "United States",
				PostedOn: "August 24, 2021 07:28 UTC",
				Text: "Description",
				Category: "Front-End Development",
				Link: "https://www.upwork.com/jobs/?source=rss",
				Skills: []string{"React Bootstrap", "Chart.js"},
				Budget: 0,
				HourlyRange: []float64{},
			},
		},{
			Src: "Description2<br /><br /><b>Hourly Range</b>: $15.00-$30.00\n\n<br /><b>Posted On</b>: August 16, 2021 12:13 UTC<br /><b>Category</b>: Front-End Development<br /><b>Skills</b>:Web Design,     Very Small (1-9 employees)\n<br /><b>Country</b>: United Kingdom\n<br /><a href=\"https://www.upwork.com/jobs/?source=rss\">click to apply</a>\n",
			Sum: Summary{
				Country: "United Kingdom",
				PostedOn: "August 16, 2021 12:13 UTC",
				Text: "Description2",
				Category: "Front-End Development",
				Link: "https://www.upwork.com/jobs/?source=rss",
				Skills: []string{"Web Design", "Very Small (1-9 employees)"},
				Budget: 0,
				HourlyRange: []float64{15.00, 30.00},
			},
		},
	}

	for i,v := range tests {
		s := ParseSummary(v.Src)
		if s.Country != v.Sum.Country {
			t.Errorf(
				"%d Country Test Error expect:%s, got %s",
				i,
				v.Sum.Country,
				s.Country,
			)
		}

		if s.PostedOn != v.Sum.PostedOn {
			t.Errorf(
				"%d PostedOn Test Error expect:%s, got %s",
				i,
				v.Sum.PostedOn,
				s.PostedOn,
			)
		}

		if s.Text != v.Sum.Text {
			t.Errorf(
				"%d Text Test Error expect:%s, got %s",
				i,
				v.Sum.Text,
				s.Text,
			)
		}

		if s.Category != v.Sum.Category {
			t.Errorf(
				"%d Category Test Error expect:%s, got %s",
				i,
				v.Sum.Category,
				s.Category,
			)
		}

		if s.Link != v.Sum.Link {
			t.Errorf(
				"%d Link Test Error expect:%s, got %s",
				i,
				v.Sum.Link,
				s.Link,
			)
		}

		if s.Budget != v.Sum.Budget {
			t.Errorf(
				"%d Budget Test Error expect:%d, got %d",
				i,
				v.Sum.Budget,
				s.Budget,
			)
		}

		if len(s.Skills) != len(v.Sum.Skills) {
			t.Errorf(
				"%d Skills Length Test Error expect:%d, got %d",
				i,
				len(v.Sum.Skills),
				len(s.Skills),
			)
		} else {
			for s,sk := range s.Skills {
				if sk != v.Sum.Skills[s] {
					t.Errorf(
						"%d Skill Test Error expect:%s, got %s",
						i,
						v.Sum.Skills[s] ,
						sk,
					)
				}
			}
		}

		if len(s.HourlyRange) != len(v.Sum.HourlyRange) {
			t.Errorf(
				"%d Hourly Range Length Test Error expect:%d, got %d",
				i,
				len(v.Sum.HourlyRange),
				len(s.HourlyRange),
			)
		} else {
			for s,sk := range s.HourlyRange {
				if sk != v.Sum.HourlyRange[s] {
					t.Errorf(
						"%d Hourly Range Test Error expect:%g, got %g",
						i,
						v.Sum.HourlyRange[s] ,
						sk,
					)
				}
			}
		}
	}
}
