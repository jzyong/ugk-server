package mode

// Server 服务器全局信息
type Server struct {
	Id     int64  `_id`    //数据库唯一ID
	RoomId uint32 `roomId` //房间ID，自行增长
	dirty  bool   //数据是否修改
}

func (s *Server) GetDirty() bool {
	return s.dirty
}

func (s *Server) SetDirty(dirty bool) {
	s.dirty = dirty
}
