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
	Feed *UWFeed
	Error error
}

func (up Upwork)Run(period time.Duration) chan *RSSResponse {
	ch := make(chan *RSSResponse)
	go func() {
		feeds := up.config.Feeds
		for {
			time.Sleep(period)
			for _, v := range feeds {
				up.FetchRss(v, ch)
			}
		}
	}()
	return ch
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
	ch <- &RSSResponse{ &rss, err }
}
