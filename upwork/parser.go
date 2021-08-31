package upwork

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Feed struct {
	Id string
	Entries []Entry
}

type Entry struct {
	Id      string
	Title   string
	Updated time.Time
	Link    string
	Summary Summary
	Content string
}

type Summary struct {
	Text string
	PostedOn string
	Category string
	Skills 	[]string
	Country string
	Link string
	Budget int
	HourlyRange []float64
}

func ConvertUWFeedToFeed(uwFeed * UWFeed) (*Feed, error) {
	res := Feed{
		Id:uwFeed.Id,
		Entries: make([]Entry, 0),
	}

	for _,entry := range uwFeed.Entries {
		e, err := ConvertUWEntryToEntry(&entry)
		if err != nil {
			return &res, err
		}

		res.Entries = append(res.Entries, *e)
	}

	return &res, nil
}


func ConvertUWEntryToEntry(uwEntry * UWEntry) (*Entry, error) {
	res := Entry{
		Id: uwEntry.Id,
		Title: uwEntry.Title,
		Link: uwEntry.Link.Href,
		Content: uwEntry.Content,
	}

	res.Summary = ParseSummary(uwEntry.Summary)

	uTime, err := time.Parse("2006-01-02T15:04:05-07:00", uwEntry.Updated)
	if err != nil {
		return &res, err
	}
	res.Updated = uTime

	return &res, nil
}

func ParseSummary(summary string) Summary {
	res := Summary{}
	brR, _ := regexp.Compile("<br />")
	splitLongR, _ := regexp.Compile("<br /><br /><br /><b>")
	splitR, _ := regexp.Compile("<br /><br /><b>")
	postedOnR, _ := regexp.Compile("<b>Posted On</b>: ")
	categoryR, _ := regexp.Compile("<b>Category</b>: ")
	skilsR, _ := regexp.Compile("<b>Skills</b>:")
	countryR, _ := regexp.Compile("<b>Country</b>: ")
	budgetR, _ := regexp.Compile("<b>Budget</b>: ")
	hourlyRangeR, _ := regexp.Compile("<b>Hourly Range</b>: ")

	startLink, _ := regexp.Compile("<br /><a href=\"")
	endLink, _ := regexp.Compile("\">click to apply</a>")


	var rest string
	descEndIndex:= splitLongR.FindStringIndex(summary)
	if len(descEndIndex) == 0 {
		descEndIndex = splitR.FindStringIndex(summary)
		if len(descEndIndex) == 0 {
			return res
		}
	}

	text := summary[:descEndIndex[0]]
	res.Text = formatSummaryText(text)
	rest = summary[descEndIndex[0]:]

	res.PostedOn = getSummaryItem(rest, postedOnR, brR)
	res.Category = getSummaryItem(rest, categoryR, brR)
	skills := getSummaryItem(rest, skilsR, brR)
	res.Skills = parseSkillsToSlice(skills)
	res.Country = getSummaryItem(rest, countryR, brR)
	res.Link = getSummaryItem(rest, startLink, endLink)
	budget := getSummaryItem(rest, budgetR, endLink)
	res.Budget = parseBudget(budget)
	hourly := getSummaryItem(rest, hourlyRangeR, brR)
	res.HourlyRange = parseHourlyRange(hourly)
	return res
}

func formatSummaryText(text string) string {
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "&amp;amp;", "&")
	text = strings.ReplaceAll(text, "&amp;gt;", ">")
	text = strings.ReplaceAll(text, "&ndash;", "-")
	text = strings.ReplaceAll(text, "&#039;", "'")
	text = strings.ReplaceAll(text, "&middot;", "·")
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&bull;", "•")
	return text
}

func getSummaryItem(text string, item, br *regexp.Regexp) string {
	res := ""
	itemIndex := item.FindStringIndex(text)
	if len(itemIndex) >= 2 {
		t := text[itemIndex[1]:]
		brIndex := br.FindStringIndex(t)
		if len(brIndex) >= 2 {
			res = t[:brIndex[0]]
		} else {
			res  = t
		}
	}

	res = strings.ReplaceAll(res, "\n", "")

	return res
}

func parseSkillsToSlice(skills string ) []string {
	skills = strings.ReplaceAll(skills, "\n", "")
	return strings.Split(skills, ",     ")
}

func parseBudget(budget string ) int{
	budget = strings.ReplaceAll(budget, "$", "")
	res, err := strconv.Atoi(budget)
	if err != nil {
		res = 0
	}
	return res
}

func parseHourlyRange(hourlyRange string ) []float64{
	hourlyRange = strings.ReplaceAll(hourlyRange, "$", "")
	hr := strings.Split(hourlyRange, "-")
	if len(hr) >=2 {
		from, errFrom := strconv.ParseFloat(hr[0], 32)
		to, errTo := strconv.ParseFloat(hr[1], 32)

		if errFrom == nil && errTo == nil {
			return []float64{from, to}
		}
	}
	return []float64{}
}
