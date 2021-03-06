package bot

import (
	"github.com/profawk/espurnaBot/api"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
	"time"
)

func apiTaskAdapter(b *tb.Bot, dest tb.Recipient, apiCall api.ApiCall) func() {
	return func() {
		sendApiMessage(b, dest, "Task executed\nThe relay is %s", apiCall)
	}
}

const newTaskMagic = `to create a new task reply to this message with the desired time in this format

"[dd/mm] hh:mm"

(day and month are optional)
to cancel simply /start`

type taskRepr struct {
	when      time.Time
	what      string
	recurring bool
}

func (t *taskRepr) UnmarshalText(text []byte) error {
	parts := strings.Split(string(text), "|")
	var tm time.Time
	err := tm.UnmarshalText([]byte(parts[0]))
	if err != nil {
		return err
	}
	t.when = tm
	t.what = parts[1]
	t.recurring, err = strconv.ParseBool(parts[2])
	return err
}

func (t taskRepr) MarshalText() (text []byte, err error) {
	tm, _ := t.when.MarshalText()
	return []byte(strings.Join([]string{
		string(tm),
		t.what,
		strconv.FormatBool(t.recurring),
	}, "|")), nil
}

func addApiHandler(b *tb.Bot, what string) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		var repr taskRepr
		repr.UnmarshalText([]byte(c.Data))
		repr.what = what
		b.EditReplyMarkup(c.Message, getAddKeyboard(repr))
		b.Respond(c)
	}
}

func addRecurringHandler(b *tb.Bot, what string) func(c *tb.Callback) {
	return func(c *tb.Callback) {
		var repr taskRepr
		repr.UnmarshalText([]byte(c.Data))
		repr.recurring = !repr.recurring
		b.EditReplyMarkup(c.Message, getAddKeyboard(repr))
		b.Respond(c)
	}
}
