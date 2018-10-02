package core


const (
	SHARD_MASTER = 0xFFFF
	SHARD_DEFAULT = 0x00
)

const (
	ST_START = 0x00
	ST_UNKNOWN = 0x01
	ST_PENDING = 0x02
	ST_DISABLED = 0x03
	ST_ENABLED = 0x04

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
		shardMask uint16
		shardPermitted []uint8

	//  Root   //common.Hash    `json:"stateRoot"       gencodec:"required"`
	  //Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
}

func NewSystemState() *SysState {
	return &SysState{
		shardId:0,
		state : ST_START,
	}
}

func (st *SysState) Reset (shardId uint16){
	st.shardId = shardId
	st.SetState(ST_PENDING)
}
func (st *SysState)SetState(newState uint8) {
	if(st.state != newState) {
		
	}

	st.state = newState;
}

func (st *SysState) GetState()  uint8 {
	return st.state
}
//check state of specified shardId
func (st *SysState) CheckStateOfShard(shardId uint16) uint8 {
	if shardId != 0x0000 {
		//first check if out of shardMask
		if ((shardId & ^ st.shardMask) != 0) { //out of ShardMask
			return (ST_DISABLED);
		}else{
			seg := shardId /8;
			offset := shardId %8;
			if (st.shardPermitted[seg] & (0x01 << offset)) != 0 {
				return (ST_ENABLED);
			}else {
				return (ST_PENDING);
			}
		}

		
	} else {  //default shard
		return ST_ENABLED;
	}
}
//change shard state according master block
func (st *SysState) OnNewMasterBlock(shardMasks uint16,shardPermitted []uint8) {
	st.shardMask = uint16(shardMasks);
	st.shardPermitted = shardPermitted;
	st.SetState(st.CheckStateOfShard(st.shardId));
	//when switched into pend from enabled, the tx pool should be retained
}
//return shard Id
func (st *SysState) Shard() uint16 {
	return st.shardId;
}

//check if a tx should be proceeded by this shard
func (st *SysState) ShouldProceed(txId uint16) bool {

	if st.shardId == SHARD_DEFAULT {
		return (txId == 0 || (st.CheckStateOfShard(txId) != ST_ENABLED))

	}else{
		return   (txId == st.shardId && st.GetState() == ST_ENABLED)
	}
}
//When this 
func (st *SysState) ShouldPend(txId uint16) bool {
	return txId != 0 && (txId == st.shardId && st.GetState() == ST_PENDING)
}
