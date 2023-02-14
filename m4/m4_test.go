package m4_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jeromelesaux/m4client/m4"
	"github.com/stretchr/testify/assert"
)

func TestNodeParsing(t *testing.T) {
	content := "Jeux/\n" +
		"Ishido.dsk,1,190K\n" +
		"Doomsday_Lost_Echoes_v1.0,0,0\n" +
		"GalacticTomb_128K,0,0\n" +
		"ImperialMahjong,0,0\n" +
		"Orion Prime (FR) (Cargosoft),0,0\n" +
		"The Shadows Of Sergoth v1.0 (F,UK,S) (128K) (Face A) (2018) [Original].dsk,1,190K\n" +
		"The Shadows Of Sergoth v1.0 (F,UK,S) (128K) (Face B) (2018) [Original].dsk,1,190K\n" +
		"Ishido,0,0"

	d := m4.NewM4Dir(content)
	if d.CurrentPath != "Jeux/" {
		t.Fatalf("Expected currentpath value Jeux/ and gets %s", d.CurrentPath)
	}
	if len(d.Nodes) != 8 {
		t.Fatalf("Expected 8 nodes and gets :%d", len(d.Nodes))
	}
}

func TestM4Client(t *testing.T) {
	t.Run("start", func(t *testing.T) {
		u := startServer(t, http.StatusOK, "")
		client := m4.M4Client{IPClient: u.Host}
		err := client.Start()
		assert.NoError(t, err)
	})
	t.Run("expectedStartFailed", func(t *testing.T) {
		u := startServer(t, http.StatusInternalServerError, "")
		client := m4.M4Client{IPClient: u.Host}
		err := client.Start()
		assert.Error(t, err)
	})
	t.Run("ResetCPC", func(t *testing.T) {
		u := startServer(t, http.StatusOK, "")
		client := m4.M4Client{IPClient: u.Host}
		err := client.ResetCpc()
		assert.NoError(t, err)
	})
	t.Run("M4Reset", func(t *testing.T) {
		u := startServer(t, http.StatusOK, "")
		client := m4.M4Client{IPClient: u.Host}
		err := client.ResetM4()
		assert.NoError(t, err)
	})
	t.Run("Pause", func(t *testing.T) {
		u := startServer(t, http.StatusOK, "")
		client := m4.M4Client{IPClient: u.Host}
		err := client.PauseCpc()
		assert.NoError(t, err)
	})
	t.Run("Mkdir", func(t *testing.T) {
		u := startServer(t, http.StatusOK, "")
		client := m4.M4Client{IPClient: u.Host}
		err := client.MakeDirectory("test")
		assert.NoError(t, err)
	})
	t.Run("Ls", func(t *testing.T) {
		r := `Jeux/
		"Ishido.dsk,1,190K
		"Doomsday_Lost_Echoes_v1.0,0,0
		"GalacticTomb_128K,0,0
		"ImperialMahjong,0,0`
		u := startServer(t, http.StatusOK, r)
		client := m4.M4Client{IPClient: u.Host}
		path, err := client.Ls("test")
		assert.NoError(t, err)
		assert.NotEmpty(t, path)
	})
	t.Run("GetDirParseFile", func(t *testing.T) {
		r := `Directory/
Ishido.dsk,1,190K`
		u := startServer(t, http.StatusOK, r)
		client := m4.M4Client{IPClient: u.Host}
		path, err := client.GetDir("test")
		assert.NoError(t, err)
		assert.NotNil(t, path)
		assert.False(t, path.Nodes[0].IsDirectory)
		assert.Equal(t, "190K", path.Nodes[0].Size)
		assert.Equal(t, "Ishido.dsk", path.Nodes[0].Name)
	})

	t.Run("GetDirParseDirectory", func(t *testing.T) {
		r := `Directory/
Subdirectory,0,0`
		u := startServer(t, http.StatusOK, r)
		client := m4.M4Client{IPClient: u.Host}
		path, err := client.GetDir("test")
		assert.NoError(t, err)
		assert.NotNil(t, path)
		assert.True(t, path.Nodes[0].IsDirectory)
		assert.Equal(t, "0", path.Nodes[0].Size)
		assert.Equal(t, "Subdirectory", path.Nodes[0].Name)
	})
}

func startServer(t *testing.T, httpCode int, httpBody string) *url.URL {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpCode)
		w.Write([]byte(httpBody))
	}))
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	return u
}
