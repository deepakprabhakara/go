// Skip this test file in Go <1.8 because it's using http.Server.Shutdown
// +build go1.8

package server

import (
	"testing"
	"time"

	"github.com/stellar/go/services/bifrost/database"
	"github.com/stellar/go/services/bifrost/litecoin"
	"github.com/stellar/go/services/bifrost/queue"
	"github.com/stellar/go/services/bifrost/sse"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LitecoinRailTestSuite struct {
	suite.Suite
	Server        *Server
	MockDatabase  *database.MockDatabase
	MockQueue     *queue.MockQueue
	MockSSEServer *sse.MockServer
}

func (suite *LitecoinRailTestSuite) SetupTest() {
	suite.MockDatabase = &database.MockDatabase{}
	suite.MockQueue = &queue.MockQueue{}
	suite.MockSSEServer = &sse.MockServer{}

	suite.Server = &Server{
		Database:          suite.MockDatabase,
		TransactionsQueue: suite.MockQueue,
		SSEServer:         suite.MockSSEServer,
		minimumValueLit:   100000000, // 1 LTC
	}
	suite.Server.initLogger()
}

func (suite *LitecoinRailTestSuite) TearDownTest() {
	suite.MockDatabase.AssertExpectations(suite.T())
	suite.MockQueue.AssertExpectations(suite.T())
	suite.MockSSEServer.AssertExpectations(suite.T())
}

func (suite *LitecoinRailTestSuite) TestInvalidValue() {
	transaction := litecoin.Transaction{
		Hash:       "109fa1c369680c2f27643fdd160620d010851a376d25b9b00ef71afe789ea6ed",
		TxOutIndex: 0,
		ValueLit:   50000000, // 0.5 LTC
		To:         "Lhd98J63jWM44tY8tcGPcvCdRDruDadyJj",
	}
	suite.MockDatabase.AssertNotCalled(suite.T(), "AddProcessedTransaction")
	suite.MockQueue.AssertNotCalled(suite.T(), "QueueAdd")
	err := suite.Server.onNewLitecoinTransaction(transaction)
	suite.Require().NoError(err)
}

func (suite *LitecoinRailTestSuite) TestAssociationNotExist() {
	transaction := litecoin.Transaction{
		Hash:       "109fa1c369680c2f27643fdd160620d010851a376d25b9b00ef71afe789ea6ed",
		TxOutIndex: 0,
		ValueLit:   100000000,
		To:         "Lhd98J63jWM44tY8tcGPcvCdRDruDadyJj",
	}
	suite.MockDatabase.
		On("GetAssociationByChainAddress", database.ChainLitecoin, "Lhd98J63jWM44tY8tcGPcvCdRDruDadyJj").
		Return(nil, nil)
	suite.MockDatabase.AssertNotCalled(suite.T(), "AddProcessedTransaction")
	suite.MockQueue.AssertNotCalled(suite.T(), "QueueAdd")
	err := suite.Server.onNewLitecoinTransaction(transaction)
	suite.Require().NoError(err)
}

func (suite *LitecoinRailTestSuite) TestAssociationAlreadyProcessed() {
	transaction := litecoin.Transaction{
		Hash:       "109fa1c369680c2f27643fdd160620d010851a376d25b9b00ef71afe789ea6ed",
		TxOutIndex: 0,
		ValueLit:   100000000,
		To:         "Lhd98J63jWM44tY8tcGPcvCdRDruDadyJj",
	}
	association := &database.AddressAssociation{
		Chain:            database.ChainLitecoin,
		AddressIndex:     1,
		Address:          "Lhd98J63jWM44tY8tcGPcvCdRDruDadyJj",
		StellarPublicKey: "GDULKYRRVOMASFMXBYD4BYFRSHAKQDREEVVP2TMH2CER3DW2KATIOASB",
		CreatedAt:        time.Now(),
	}
	suite.MockDatabase.
		On("GetAssociationByChainAddress", database.ChainLitecoin, transaction.To).
		Return(association, nil)
	suite.MockDatabase.
		On("AddProcessedTransaction", database.ChainLitecoin, transaction.Hash, transaction.To).
		Return(true, nil)
	suite.MockQueue.AssertNotCalled(suite.T(), "QueueAdd")
	err := suite.Server.onNewLitecoinTransaction(transaction)
	suite.Require().NoError(err)
}

func (suite *LitecoinRailTestSuite) TestAssociationSuccess() {
	transaction := litecoin.Transaction{
		Hash:       "109fa1c369680c2f27643fdd160620d010851a376d25b9b00ef71afe789ea6ed",
		TxOutIndex: 0,
		ValueLit:   100000000,
		To:         "Lhd98J63jWM44tY8tcGPcvCdRDruDadyJj",
	}
	association := &database.AddressAssociation{
		Chain:            database.ChainLitecoin,
		AddressIndex:     1,
		Address:          "Lhd98J63jWM44tY8tcGPcvCdRDruDadyJj",
		StellarPublicKey: "GDULKYRRVOMASFMXBYD4BYFRSHAKQDREEVVP2TMH2CER3DW2KATIOASB",
		CreatedAt:        time.Now(),
	}
	suite.MockDatabase.
		On("GetAssociationByChainAddress", database.ChainLitecoin, transaction.To).
		Return(association, nil)
	suite.MockDatabase.
		On("AddProcessedTransaction", database.ChainLitecoin, transaction.Hash, transaction.To).
		Return(false, nil)
	suite.MockQueue.
		On("QueueAdd", mock.AnythingOfType("queue.Transaction")).
		Return(nil).
		Run(func(args mock.Arguments) {
			queueTransaction := args.Get(0).(queue.Transaction)
			suite.Assert().Equal(transaction.Hash, queueTransaction.TransactionID)
			suite.Assert().Equal("LTC", string(queue.AssetCodeLTC))
			suite.Assert().Equal(queue.AssetCodeLTC, queueTransaction.AssetCode)
			suite.Assert().Equal("1.0000000", queueTransaction.Amount)
			suite.Assert().Equal(association.StellarPublicKey, queueTransaction.StellarPublicKey)
		})
	suite.MockSSEServer.
		On("BroadcastEvent", transaction.To, sse.TransactionReceivedAddressEvent, []byte(nil))
	err := suite.Server.onNewLitecoinTransaction(transaction)
	suite.Require().NoError(err)
}

func TestLitecoinRailTestSuite(t *testing.T) {
	suite.Run(t, new(LitecoinRailTestSuite))
}
