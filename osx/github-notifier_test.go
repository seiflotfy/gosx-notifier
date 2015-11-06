package githubnotifier

import (
	"log"
	"path/filepath"
	"testing"
)

func TestInstall(t *testing.T) {
	//assert file exists

	if !exists(finalPath) {
		t.Error("TestInstall failed to install the terminal-notifier.app bundle")
	} else {
		log.Println("terminal-notifier.app bundle installed successfully at: ", finalPath)
	}
}

func TestNewNotifier(t *testing.T) {
	n := NewNotification("Hello")

	//assert defaults
	if n.Message != "Hello" {
		t.Error("NewNotification doesn't have a Message specified")
	}
}

func TestPush(t *testing.T) {
	n := NewNotification("Testing Push")
	err := n.Push()

	if err != nil {
		t.Error("TestPush failed with error: ", err)
	}
}

func TestTitle(t *testing.T) {
	n := NewNotification("Testing Title")
	n.Title = "gosx-notifier is amazing!"
	err := n.Push()

	if err != nil {
		t.Error("TestTitle failed with error: ", err)
	}
}

func TestSubtitle(t *testing.T) {
	n := NewNotification("Testing Subtitle")
	n.Subtitle = "gosx-notifier rocks!"

	err := n.Push()

	if err != nil {
		t.Error("TestSubtitle failed with error: ", err)
	}
}

func TestSender(t *testing.T) {

	for _, s := range []string{"com.apple.Safari", "com.apple.iTunes"} {

		n := NewNotification("Testing Icon")
		n.Title = s
		n.Sender = s

		err := n.Push()

		if err != nil {
			t.Error("TestSender failed with error: ", err)
		}
	}
}

func TestGroup(t *testing.T) {
	const appID string = "github.com/deckarep/gosx-notifier"

	for i := 0; i < 3; i++ {
		n := NewNotification("Testing Group Functionality...")
		n.Group = appID

		err := n.Push()

		if err != nil {
			t.Error("TestGroup failed with error: ", err)
		}

	}
}

func TestAppIcon(t *testing.T) {
	const appIcon string = "gopher.png"

	n := NewNotification("Testing App Icon")

	if icon, err := filepath.Abs(appIcon); err != nil {
		t.Error("TestAppIcon could not get the absolute file of: ", appIcon)
	} else {
		n.AppIcon = icon
	}

	err := n.Push()

	if err != nil {
		t.Error("TestAppIcon failed with error: ", err)
	}
}

func TestContentImage(t *testing.T) {
	const contentImage string = "gopher.png"

	n := NewNotification("Testing Content Image")

	if img, err := filepath.Abs(contentImage); err != nil {
		t.Error("TestAppIcon could not get the absolute file of: ", contentImage)
	} else {
		n.ContentImage = img
	}

	err := n.Push()

	if err != nil {
		t.Error("TestContentImage failed with error: ", err)
	}
}

func TestContentImageAndIcon(t *testing.T) {
	const image string = "gopher.png"

	n := NewNotification("Testing Content Image and Icon")
	n.Title = "Hey Gopher!"
	n.Subtitle = "I eat Goroutines for breakfast!"

	if img, err := filepath.Abs(image); err != nil {
		t.Error("TestAppIcon could not get the absolute file of: ", image)
	} else {
		n.ContentImage = img
		n.AppIcon = img
	}

	err := n.Push()

	if err != nil {
		t.Error("TestContentImageAndIcon failed with error: ", err)
	}
}

/*
	Not an easy way to verify the tests below actually work as designed, but here for completion.
*/

func TestSound(t *testing.T) {
	n := NewNotification("Testing Sound")
	n.Sound = Default
	err := n.Push()

	if err != nil {
		t.Error("TestSound failed with error: ", err)
	}
}

func TestLink_Url(t *testing.T) {
	n := NewNotification("Testing Link Url")
	n.Link = "http://www.yahoo.com"
	err := n.Push()

	if err != nil {
		t.Error("TestLink failed with error: ", err)
	}
}

func TestLink_App_Bundle(t *testing.T) {
	n := NewNotification("Testing Link Terminal")
	n.Link = "com.apple.Safari"
	err := n.Push()

	if err != nil {
		t.Error("TestLink failed with error: ", err)
	}
}
