package provider

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/hirochachacha/go-smb2"
)

type SMBInput struct {
	Addr     string
	User     string
	Password string
	Domain   string
}

type SMBProvider struct {
	filepath string
	addr     string
	user     string
	password string
	domain   string

	session *smb2.Session
	share   *smb2.Share
}

// addr - ip:port
func NewSMBManager(input SMBInput) (Provider, error) {
	if input.Addr == "" {
		return nil, fmt.Errorf("addr smb is empty")
	}

	smb := &SMBProvider{
		addr:     input.Addr,
		user:     input.User,
		password: input.Password,
		domain:   input.Domain,
	}

	host, share, path, err := smb.parseUNC(input.Addr)
	if err != nil {
		return nil, err
	}

	session, err := smb.startSession(host)
	if err != nil {
		return nil, err
	}

	if err := smb.connectToShare(session, share); err != nil {
		return nil, err
	}

	smb.filepath = path

	return smb, nil
}

func (m *SMBProvider) parseUNC(path string) (host, share, filePath string, err error) {
	path = strings.TrimPrefix(path, `\\`)

	parts := strings.Split(path, `\`)
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid UNC path")
	}

	host = parts[0]
	share = parts[1]
	filePath = strings.Join(parts[2:], `\`)

	return
}

func (m *SMBProvider) connectToShare(s *smb2.Session, share string) error {
	fs, err := s.Mount(share)
	if err != nil {
		return err
	}

	m.share = fs
	return nil
}

func (m *SMBProvider) Get() ([]byte, error) {
	src, err := m.share.Open(m.filepath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	var data []byte
	data, err = io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *SMBProvider) Disconnect() {
	m.session.Logoff()
	m.share.Umount()
}

func (m *SMBProvider) startSession(host string) (*smb2.Session, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// auth
	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     m.user,
			Password: m.password,
			Domain:   m.domain, // или WORKGROUP / DOMAIN
		},
	}

	// session
	s, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}

	return s, nil
}
