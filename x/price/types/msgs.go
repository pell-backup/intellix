package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RouterKey                     = "price"
	TypeVoteFinalizedRequestPrice = "vote_finalized_request_price"
)

var (
	amino    = codec.NewLegacyAmino()
	AminoCdc = codec.NewAminoCodec(amino)
)

// Route returns the route of MsgEditDataSource - "oracle" (sdk.Msg interface).
func (m MsgVoteFinalizedRequestPrice) Route() string { return RouterKey }

// Type returns the message type of MsgEditDataSource (sdk.Msg interface).
func (m MsgVoteFinalizedRequestPrice) Type() string { return TypeVoteFinalizedRequestPrice }

// ValidateBasic checks whether the given MsgEditDataSource instance (sdk.Msg interface).
func (m MsgVoteFinalizedRequestPrice) ValidateBasic() error {
	if m.ValidatedData == nil || m.PriceFeedResponse == nil {
		return fmt.Errorf("empty data")
	}
	return nil
}

// GetSigners returns the required signers for the given MsgEditDataSource (sdk.Msg interface).
func (m MsgVoteFinalizedRequestPrice) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.TaskRaw.CallbackAddress)
	return []sdk.AccAddress{sender}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (m MsgVoteFinalizedRequestPrice) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&m))
}

// Route returns the route of MsgEditDataSource - "oracle" (sdk.Msg interface).
func (m MsgVoteRequestPriceFeed) Route() string { return RouterKey }

// Type returns the message type of MsgEditDataSource (sdk.Msg interface).
func (m MsgVoteRequestPriceFeed) Type() string { return TypeVoteFinalizedRequestPrice }

// ValidateBasic checks whether the given MsgEditDataSource instance (sdk.Msg interface).
func (m MsgVoteRequestPriceFeed) ValidateBasic() error {
	if m.RequestId == nil {
		return fmt.Errorf("empty data")
	}
	return nil
}

// GetSigners returns the required signers for the given MsgEditDataSource (sdk.Msg interface).
func (m MsgVoteRequestPriceFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

// GetSignBytes returns raw JSON bytes to be signed by the signers (sdk.Msg interface).
func (m MsgVoteRequestPriceFeed) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(&m))
}
