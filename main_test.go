package main

import (
	"github.com/RocketChat/Rocket.Chat.Go.SDK/rest"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormattedEventAction(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")

	action := formattedEventAction(event)
	assert.Equal("RESOLVED", action)

	event.Check.Status = 1
	action = formattedEventAction(event)
	assert.Equal("ALERT", action)
}

func TestEventKey(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	eventKey := eventKey(event)
	assert.Equal("entity1/check1", eventKey)
}

func TestEventSummary(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	event.Check.Output = "disk is full"

	eventKey := eventSummary(event, 100)
	assert.Equal("entity1/check1:disk is full", eventKey)

	eventKey = eventSummary(event, 5)
	assert.Equal("entity1/check1:disk ...", eventKey)
}

func TestFormattedMessage(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	event.Check.Output = "disk is full"
	event.Check.Status = 1
	formattedMsg := formattedMessage(event)
	assert.Equal("ALERT - entity1/check1:disk is full", formattedMsg)
}

func TestMessageColor(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")

	event.Check.Status = 0
	color := messageColor(event)
	assert.Equal("good", color)

	event.Check.Status = 1
	color = messageColor(event)
	assert.Equal("warning", color)

	event.Check.Status = 2
	color = messageColor(event)
	assert.Equal("danger", color)
}

func TestMessageStatus(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")

	event.Check.Status = 0
	status := messageStatus(event)
	assert.Equal("Resolved", status)

	event.Check.Status = 1
	status = messageStatus(event)
	assert.Equal("Warning", status)

	event.Check.Status = 2
	status = messageStatus(event)
	assert.Equal("Critical", status)
}

//func TestSendMessage(t *testing.T) {
//	assert := assert.New(t)
//	event := corev2.FixtureEvent("entity1", "check1")
//
//	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		body, _ := ioutil.ReadAll(r.Body)
//		expectedBody := `{"channel":"#test","attachments":[{"color":"good","fallback":"RESOLVED - entity1/check1:","title":"Description","text":"","fields":[{"title":"Status","value":"Resolved","short":false},{"title":"Entity","value":"entity1","short":true},{"title":"Check","value":"check1","short":true}]}]}`
//		assert.Equal(expectedBody, string(body))
//		w.WriteHeader(http.StatusOK)
//		_, err := w.Write([]byte(`{"ok": true}`))
//		require.NoError(t, err)
//	}))
//
//	config.rocketchatUrl = apiStub.URL
//	config.rocketchatUsername = "sensu"
//	config.rocketchatPassword = "sensu"
//	config.rocketchatChannel = "sandbox"
//	config.rocketchatTemplate = defaultTemplate
//
//	err := sendMessage(event)
//	assert.NoError(err)
//}

//func TestSendRealMessage(t *testing.T) {
//	assert := assert.New(t)
//	event := corev2.FixtureEvent("entity1", "check1")
//
//	//var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	//	body, _ := ioutil.ReadAll(r.Body)
//	//	expectedBody := `{"channel":"#test","attachments":[{"color":"good","fallback":"RESOLVED - entity1/check1:","title":"Description","text":"","fields":[{"title":"Status","value":"Resolved","short":false},{"title":"Entity","value":"entity1","short":true},{"title":"Check","value":"check1","short":true}]}]}`
//	//	assert.Equal(expectedBody, string(body))
//	//	w.WriteHeader(http.StatusOK)
//	//	_, err := w.Write([]byte(`{"ok": true}`))
//	//	require.NoError(t, err)
//	//}))
//
//	config.rocketchatUrl = ...
//	config.rocketchatUsername = ...
//	config.rocketchatPassword = ...
//	config.rocketchatChannel = ...
//
//	config.rocketchatTemplate = defaultTemplate
//
//	err := sendMessage(event)
//
//	assert.NoError(err)
//}

//func TestMainMethod(t *testing.T) {
//	assert := assert.New(t)
//	file, _ := ioutil.TempFile(os.TempDir(), "sensu-handler-rocketchat-")
//	defer func() {
//		_ = os.Remove(file.Name())
//	}()
//
//	event := corev2.FixtureEvent("entity1", "check1")
//	eventJSON, _ := json.Marshal(event)
//	_, err := file.WriteString(string(eventJSON))
//	require.NoError(t, err)
//	require.NoError(t, file.Sync())
//	_, err = file.Seek(0, 0)
//	require.NoError(t, err)
//	os.Stdin = file
//	requestReceived := false
//
//	var apiStub = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		requestReceived = true
//		w.WriteHeader(http.StatusOK)
//		_, err := w.Write([]byte(`{"ok": true}`))
//		require.NoError(t, err)
//	}))
//
//	oldArgs := os.Args
//	os.Args = []string{"rocketchat", "-w", apiStub.URL}
//	defer func() { os.Args = oldArgs }()
//
//	main()
//	assert.True(requestReceived)
//}

func TestRocket_GetServerInfo(t *testing.T) {
	rocket := rest.Client{Protocol: "https", Host: "open.rocket.chat", Port: "443"}

	info, err := rocket.GetServerInfo()

	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.NotEmpty(t, info.Version)
}
