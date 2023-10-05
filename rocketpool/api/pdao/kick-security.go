package pdao

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/types/api"
	"github.com/rocket-pool/smartnode/shared/utils/eth1"
	"github.com/urfave/cli"
)

func canProposeKickFromSecurityCouncil(c *cli.Context, address common.Address) (*api.PDAOCanProposeKickFromSecurityCouncilResponse, error) {
	// Get services
	if err := services.RequireNodeTrusted(c); err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	cfg, err := services.GetConfig(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}
	bc, err := services.GetBeaconClient(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.PDAOCanProposeKickFromSecurityCouncilResponse{}

	// Get node account
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}

	// Try proposing
	message := fmt.Sprintf("kick %s from the security council", address.Hex())
	blockNumber, pollard, encodedPollard, err := createPollard(rp, cfg, bc)
	if err != nil {
		return nil, err
	}
	gasInfo, err := protocol.EstimateProposeKickFromSecurityCouncilGas(rp, message, address, blockNumber, pollard, opts)
	if err != nil {
		return nil, err
	}

	// Update & return response
	response.BlockNumber = blockNumber
	response.Pollard = encodedPollard
	response.GasInfo = gasInfo
	return &response, nil
}

func proposeKickFromSecurityCouncil(c *cli.Context, address common.Address, blockNumber uint32, pollard string) (*api.PDAOProposeKickFromSecurityCouncilResponse, error) {
	// Get services
	if err := services.RequireNodeTrusted(c); err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Response
	response := api.PDAOProposeKickFromSecurityCouncilResponse{}

	// Get node account
	opts, err := w.GetNodeAccountTransactor()
	if err != nil {
		return nil, err
	}

	// Override the provided pending TX if requested
	err = eth1.CheckForNonceOverride(c, opts)
	if err != nil {
		return nil, fmt.Errorf("Error checking for nonce override: %w", err)
	}

	// Propose
	message := fmt.Sprintf("kick %s from the security council", address.Hex())
	truePollard, err := decodePollard(pollard)
	if err != nil {
		return nil, err
	}
	proposalID, hash, err := protocol.ProposeKickFromSecurityCouncil(rp, message, address, blockNumber, truePollard, opts)
	if err != nil {
		return nil, err
	}

	// Update & return response
	response.ProposalId = proposalID
	response.TxHash = hash
	return &response, nil
}
