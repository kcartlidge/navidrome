package api_test

import (
	"testing"

	"github.com/deluan/gosonic/api/responses"
	"github.com/deluan/gosonic/domain"
	. "github.com/deluan/gosonic/tests"
	"github.com/deluan/gosonic/tests/mocks"
	"github.com/deluan/gosonic/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetMusicDirectory(t *testing.T) {
	Init(t, false)

	mockArtistRepo := mocks.CreateMockArtistRepo()
	utils.DefineSingleton(new(domain.ArtistRepository), func() domain.ArtistRepository {
		return mockArtistRepo
	})
	mockAlbumRepo := mocks.CreateMockAlbumRepo()
	utils.DefineSingleton(new(domain.AlbumRepository), func() domain.AlbumRepository {
		return mockAlbumRepo
	})
	mockMediaFileRepo := mocks.CreateMockMediaFileRepo()
	utils.DefineSingleton(new(domain.MediaFileRepository), func() domain.MediaFileRepository {
		return mockMediaFileRepo
	})

	Convey("Subject: GetMusicDirectory Endpoint", t, func() {
		Convey("Should fail if missing Id parameter", func() {
			_, w := Get(AddParams("/rest/getMusicDirectory.view"), "TestGetMusicDirectory")

			So(w.Body, ShouldReceiveError, responses.ERROR_MISSING_PARAMETER)
		})
		Convey("Id is for an artist", func() {
			Convey("Return fail on Artist Table error", func() {
				mockArtistRepo.SetError(true)
				_, w := Get(AddParams("/rest/getMusicDirectory.view", "id=1"), "TestGetMusicDirectory")

				So(w.Body, ShouldReceiveError, responses.ERROR_GENERIC)
			})
		})
		Convey("When id is not found", func() {
			mockArtistRepo.SetData(`[{"Id":"1","Name":"The Charlatans"}]`, 1)
			_, w := Get(AddParams("/rest/getMusicDirectory.view", "id=NOT_FOUND"), "TestGetMusicDirectory")

			So(w.Body, ShouldReceiveError, responses.ERROR_DATA_NOT_FOUND)
		})
		Convey("When id matches an artist", func() {
			mockArtistRepo.SetData(`[{"Id":"1","Name":"The KLF"}]`, 1)

			Convey("Without albums", func() {
				_, w := Get(AddParams("/rest/getMusicDirectory.view", "id=1"), "TestGetMusicDirectory")

				So(w.Body, ShouldContainJSON, `"id":"1","name":"The KLF"`)
			})
			Convey("With albums", func() {
				mockAlbumRepo.SetData(`[{"Id":"A","Name":"Tardis","ArtistId":"1"}]`, 1)
				_, w := Get(AddParams("/rest/getMusicDirectory.view", "id=1"), "TestGetMusicDirectory")

				So(w.Body, ShouldContainJSON, `"child":[{"album":"Tardis","artist":"The KLF","id":"A","isDir":true,"title":"Tardis"}]`)
			})
		})
		Convey("When id matches an album with tracks", func() {
			mockArtistRepo.SetData(`[{"Id":"2","Name":"Céu"}]`, 1)
			mockAlbumRepo.SetData(`[{"Id":"A","Name":"Vagarosa","ArtistId":"2"}]`, 1)
			mockMediaFileRepo.SetData(`[{"Id":"3","Title":"Cangote","AlbumId":"A"}]`, 1)
			_, w := Get(AddParams("/rest/getMusicDirectory.view", "id=A"), "TestGetMusicDirectory")

			So(w.Body, ShouldContainJSON, `"child":[{"id":"3","isDir":false,"title":"Cangote"}]`)
		})
		Reset(func() {
			mockArtistRepo.SetData("[]", 0)
			mockArtistRepo.SetError(false)

			mockAlbumRepo.SetData("[]", 0)
			mockAlbumRepo.SetError(false)

			mockMediaFileRepo.SetData("[]", 0)
			mockMediaFileRepo.SetError(false)
		})
	})
}
