package integration_test

import (
	cosmosMath "cosmossdk.io/math"
	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/stretchr/testify/require"
)

// register alice as a reputer in topic 1, then check success
func StakeAliceAsReputerTopic1(m TestMetadata) {
	addStake := &emissionstypes.MsgAddStake{
		Sender:  m.n.AliceAddr,
		TopicId: 1,
		Amount:  cosmosMath.NewUint(1000000),
	}

	txResp, err := m.n.Client.BroadcastTx(m.ctx, m.n.AliceAcc, addStake)
	require.NoError(m.t, err)
	_, err = m.n.Client.WaitForTx(m.ctx, txResp.TxHash)
	require.NoError(m.t, err)

	// Check Alice has stake on the topic
	aliceStaked, err := m.n.QueryEmissions.GetReputerStakeInTopic(
		m.ctx,
		&emissionstypes.QueryReputerStakeInTopicRequest{
			TopicId: 1,
			Address: m.n.AliceAddr,
		},
	)
	require.NoError(m.t, err)
	require.True(m.t, aliceStaked.Amount.Equal(cosmosMath.NewUint(1000000)))
}

func CheckTopic1Activated(m TestMetadata) {
	// Fetch only active topics
	pagi := &emissionstypes.QueryActiveTopicsRequest{
		Pagination: &emissionstypes.SimpleCursorPaginationRequest{
			Limit: 10,
		},
	}
	activeTopics, err := m.n.QueryEmissions.GetActiveTopics(
		m.ctx,
		pagi)
	require.NoError(m.t, err, "Fetching active topics should not produce an error")

	// Verify the correct number of active topics is retrieved
	require.Equal(m.t, len(activeTopics.Topics), 1, "Should retrieve exactly one active topic")
}

// Register two actors and check their registrations went through
func StakingChecks(m TestMetadata) {
	m.t.Log("--- Staking Alice as Reputer ---")
	StakeAliceAsReputerTopic1(m)
	m.t.Log("--- Check reactivating Topic 1 ---")
	CheckTopic1Activated(m)
}
