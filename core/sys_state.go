package core



const (
	ST_START = 0x00
	ST_UNKNOWN = 0x01
	ST_DISABLED = 0x02
	ST_ENABLED = 0x03
)
/**func
var enterFuncMap = map[uint8] interface{}{
	ST_START:      func,
	ST_UNKNOWN:                  "Invalid message",
	ST_DISABLED:          "Invalid message code",
	ST_ENABLED: "Protocol version mismatch",

}
*/
type   SysState struct{
		shardId   uint16
	  state   uint8 
	//  Root   //common.Hash    `json:"stateRoot"       gencodec:"required"`
	  //Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
}

func NewSystemState() *SysState {
	return &SysState{
		shardId:0,
		state : ST_START,
	}
}

func (st SysState) Reset (shardId uint16){
	st.shardId = shardId
	st.SetState(ST_UNKNOWN)
}
func (st SysState)SetState(newState uint8) {
	if(st.state != newState) {
		
	}

	st.state = newState;
}

func (st SysState) GetState()  uint8 {
	return st.state
}

func (st SysState) OnNewBlockInfo(shardMasks uint16,shardPermitted []uint8) {

}

