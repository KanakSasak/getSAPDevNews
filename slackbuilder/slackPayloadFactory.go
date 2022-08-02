package slackbuilder

import (
	"encoding/json"
	"fmt"
	"getSAPDevNews/model"
	"github.com/slack-go/slack"
	"time"
)

func Build(data []model.News) string {

	return sapNews(data)
}

func sapNews(data []model.News) string {

	now := time.Now()
	headerText := slack.NewTextBlockObject("plain_text", ":newspaper:  SAP Blog Daily Newsletter  :newspaper:", false, false)
	headerSection := slack.NewHeaderBlock(headerText)
	contextText := slack.NewTextBlockObject("mrkdwn", "*"+now.Format("January 02, 2006")+"*", false, false)
	contextSection := slack.NewContextBlock("", []slack.MixedElement{contextText}...)
	dev := slack.NewDividerBlock()

	msg := slack.NewBlockMessage(
		headerSection,
		contextSection,
		dev,
	)
	for _, dat := range data {

		titleText := slack.NewTextBlockObject("mrkdwn", "*"+dat.Title+"*", false, false)
		titleSection := slack.NewSectionBlock(titleText, nil, nil)
		contextsubtitleText := slack.NewTextBlockObject("mrkdwn", "*"+dat.Date+"*  | *"+dat.Author+"* | "+dat.Category+"", false, false)
		contextsubtitleSection := slack.NewContextBlock("", []slack.MixedElement{contextsubtitleText}...)

		// Link Buttons
		readmoreBtnTxt := slack.NewTextBlockObject("plain_text", "Read More", false, false)
		readmoreBtn := slack.NewButtonBlockElement("", dat.Link, readmoreBtnTxt)
		mainText := slack.NewTextBlockObject("mrkdwn", dat.Desc, false, false)
		mainSection := slack.NewSectionBlock(mainText, nil, slack.NewAccessory(readmoreBtn))

		msg = slack.AddBlockMessage(msg, titleSection)
		msg = slack.AddBlockMessage(msg, contextsubtitleSection)
		msg = slack.AddBlockMessage(msg, mainSection)
		msg = slack.AddBlockMessage(msg, dev)
	}

	b, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		fmt.Println(err)
		panic(err)
		return ""
	}

	//fmt.Println(string(b))

	return string(b)
}
