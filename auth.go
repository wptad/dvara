package dvara

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"gopkg.in/mgo.v2-unstable/bson"
)

// Credential holds details to authenticate with a MongoDB server.
type Credential struct {
	// Username and Password hold the basic details for authentication.
	// Password is optional with some authentication mechanisms.
	Username string
	Password string

	// Source is the database used to establish credentials and privileges
	// with a MongoDB server. Defaults to the default database provided
	// during dial, or "admin" if that was unset.
	Source string
}

type authCmd struct {
	Authenticate int

	Nonce string
	User  string
	Key   string
}

type authResult struct {
	ErrMsg string
	Ok     bool
}

type getNonceCmd struct {
	GetNonce int
}

type getNonceResult struct {
	Nonce string
	Err   string "$err"
	Code  int
}

func (socket *mongoSocket) getNonce() (nonce string, err error) {
	fmt.Printf("Socket %p to %s: requesting a new nonce\n", socket, socket.addr)
	op := &queryOp{}
	op.query = &getNonceCmd{GetNonce: 1}
	op.collection = "admin.$cmd"
	op.limit = -1
	op.replyFunc = func(err error, reply *replyOp, docNum int, docData []byte) {
		if err != nil {
			socket.kill(errors.New("getNonce: "+err.Error()), true)
			return
		}
		result := &getNonceResult{}
		err = bson.Unmarshal(docData, &result)
		if err != nil {
			socket.kill(errors.New("Failed to unmarshal nonce: "+err.Error()), true)
			return
		}
		fmt.Println("Socket %p to %s: nonce unmarshalled: %#v", socket, socket.addr, result)
		if result.Code == 13390 {
			// mongos doesn't yet support auth (see http://j.mp/mongos-auth)
			result.Nonce = "mongos"
		} else if result.Nonce == "" {
			var msg string
			if result.Err != "" {
				msg = fmt.Sprintf("Got an empty nonce: %s (%d)", result.Err, result.Code)
			} else {
				msg = "Got an empty nonce"
			}
			socket.kill(errors.New(msg), true)
			return
		}
		nonce = result.Nonce
	}

	err = socket.Query(op)
	if err != nil {
		socket.kill(errors.New("resetNonce: "+err.Error()), true)
	}
	return
}

func (socket *mongoSocket) Login(cred Credential) error {
	nonce, err := socket.getNonce()
	if err != nil {
		return err
	}

	psum := md5.New()
	psum.Write([]byte(cred.Username + ":mongo:" + cred.Password))

	ksum := md5.New()
	ksum.Write([]byte(nonce + cred.Username))
	ksum.Write([]byte(hex.EncodeToString(psum.Sum(nil))))

	key := hex.EncodeToString(ksum.Sum(nil))

	cmd := authCmd{Authenticate: 1, User: cred.Username, Nonce: nonce, Key: key}
	res := authResult{}
	op := queryOp{}
	op.query = &cmd
	op.collection = cred.Source + ".$cmd"
	op.limit = -1
	op.replyFunc = func(err error, reply *replyOp, docNum int, docData []byte) {
		err = bson.Unmarshal(docData, &res)
		if err != nil {
			return
		}
	}
	err = socket.Query(&op)
	if err != nil {
		return err
	}
	return err
}
