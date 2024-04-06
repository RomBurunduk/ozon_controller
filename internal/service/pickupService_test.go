package service

import (
	"log"
	"os"
	"testing"

	"pvz_controller/internal/model"
)

func Test_AddPVZ(t *testing.T) {
	t.Parallel()
	var (
		c = make(chan string)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		//arrange
		s := setStorageUp(t)
		defer s.tearDown()
		s.mockPickup.EXPECT().WritePVZ(model.Pickups{
			Name:    "some",
			Address: "some",
			Contact: "some",
		}).Return(nil)

		input := []byte("some\nsome\nsome\n")
		tempFile, err := os.CreateTemp("", "example")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tempFile.Name())

		if _, err = tempFile.Write(input); err != nil {
			log.Fatal(err)
		}

		if _, err = tempFile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}

		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		os.Stdin = tempFile
		//act
		go s.srv.AddPVZ(c)
		<-c
		<-c
		//assert
	})
}
