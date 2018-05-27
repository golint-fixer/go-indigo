package discovery

import "indogo/src/networking"

type nodeDatabase struct {
	lvl     *leveldb.DB
	selfRef networking.NodeID
}

func newNodeDatabase(selfRef networking.NodeID) *nodeDatabase {

}
