package sessionMaker

import (
	"encoding/json"
	"fmt"

	"github.com/gotd/td/session"
	"github.com/jaskaur18/gotgproto/functions"
	"github.com/jaskaur18/gotgproto/storage"
)

// SessionName object consists of name and SessionType.
type SessionName struct {
	name        string
	sessionType SessionType
	data        []byte
	err         error
}

// SessionType is the type of session you want to log in through.
// It consists of three types: Session, StringSession, TelethonSession.
type SessionType int

const (
	// Session should be used for authorizing into session with default settings.
	Session SessionType = iota
	// StringSession is used as SessionType when you want to log in through the string session made by gotgproto.
	StringSession
	// TelethonSession is used as SessionType when you want to log in through the string session made by telethon - a Python MTProto library.
	TelethonSession
	// PyrogramSession is used as SessionType when you want to log in through the string session made by pyrogram - a Python MTProto library.
	PyrogramSession
)

// NewSession creates a new session with provided name string and SessionType.
func NewSession(sessionName string, sessionType SessionType) *SessionName {
	s := SessionName{
		name:        sessionName,
		sessionType: sessionType,
	}
	s.data, s.err = s.load()
	return &s
}

func (s *SessionName) load() ([]byte, error) {
	switch s.sessionType {
	case PyrogramSession:
		storage.Load("pyrogram.session", false)
		sd, err := DecodePyrogramSession(s.name)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(jsonData{
			Version: storage.LatestVersion,
			Data:    *sd,
		})
		return data, err
	case TelethonSession:
		storage.Load("telethon.session", false)
		sd, err := session.TelethonSession(s.name)
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(jsonData{
			Version: storage.LatestVersion,
			Data:    *sd,
		})
		return data, err
	case StringSession:
		storage.Load("gotgproto.session", false)
		sd, err := functions.DecodeStringToSession(s.name)
		if err != nil {
			return nil, err
		}

		// data, err := json.Marshal(jsonData{
		// 	Version: latestVersion,
		// 	Data:    *sd,
		// })
		return sd.Data, err
	default:
		if s.name == "" {
			s.name = "new"
		}
		storage.Load(fmt.Sprintf("%s.session", s.name), false)
		sFD := storage.GetSession()
		return sFD.Data, nil
	}
}

// GetName is used for retrieving name of the session.
func (s *SessionName) GetName() string {
	return s.name
}

// GetData is used for retrieving session data through provided SessionName type.
func (s *SessionName) GetData() ([]byte, error) {
	return s.data, s.err
}
