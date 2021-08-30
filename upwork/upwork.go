package upwork

import (
	"encoding/xml"
	"fmt"
	"github.com/flasherup/telegrum_upwork_bot/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type Upwork struct {
	config utils.UpworkCfg
}

func NewUpwork(config utils.UpworkCfg) *Upwork {
	return &Upwork{
		config,
	}
}

type RSSResponse struct {
	Feed *Feed
	Error error
}

func (up Upwork)Run(period time.Duration) chan *RSSResponse {
	ch := make(chan *RSSResponse)
	go func() {
		up.UpdateFeeds(ch)// Immediately
		for {
			time.Sleep(period)
			up.UpdateFeeds(ch)
		}
	}()
	return ch
}

func (up Upwork) UpdateFeeds(ch chan *RSSResponse)  {
	for _, v := range up.config.Feeds {
		up.FetchRss(v, ch)
	}
}

func (up Upwork)FetchRss(topic string, ch chan *RSSResponse) {
	url := fmt.Sprintf("%s?securityToken=%s&userUid=%s&orgUid=%s&topic=%s",
			up.config.Url,
			up.config.SecurityToken,
			up.config.UserUid,
			up.config.OrgUid,
			topic)


	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- &RSSResponse{ nil, err }
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		ch <- &RSSResponse{ nil, err }
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- &RSSResponse{ nil, err }
		return
	}

	rss := UWFeed{}
	xml.Unmarshal(content, &rss)
	res, err := ConvertUWFeedToFeed(&rss)
	ch <- &RSSResponse{ res, err }
}
