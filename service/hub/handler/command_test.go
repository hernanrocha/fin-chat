package handler

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/hernanrocha/fin-chat/messenger"
	"github.com/hernanrocha/fin-chat/service/hub/mocks"
	"github.com/hernanrocha/fin-chat/service/viewmodels"
)

type CommandMessageHandlerSuite struct {
	suite.Suite
	mockDb        sqlmock.Sqlmock
	DB            *gorm.DB
	mockHub       *mocks.MockHub
	mockMessenger *MockBotCommandMessenger
}

func TestCommandMessageHandlerSuite(t *testing.T) {
	suite.Run(t, new(CommandMessageHandlerSuite))
}

func (suite *CommandMessageHandlerSuite) SetupTest() {
	sqlDb, mock, err := sqlmock.New()
	require.NoError(suite.T(), err)

	DB, err := gorm.Open("postgres", sqlDb)
	require.NoError(suite.T(), err)

	suite.mockDb = mock
	suite.DB = DB

	suite.mockHub = mocks.NewMockHub()
	suite.mockMessenger = NewMockBotCommandMessenger()
}

func (suite *CommandMessageHandlerSuite) TestGetID() {
	expectTestUser(suite.mockDb)

	ID := "random-id"
	handler, err := NewCmdMessageHandler(ID, suite.mockMessenger, suite.mockHub, suite.DB)
	require.Nil(suite.T(), err)
	assert.Equal(suite.T(), ID, handler.GetID())

	assert.NoError(suite.T(), suite.mockDb.ExpectationsWereMet())
	suite.mockHub.AssertExpectations(suite.T())
	suite.mockMessenger.AssertExpectations(suite.T())
}

func (suite *CommandMessageHandlerSuite) TestHandleMessageNotCommand() {
	expectTestUser(suite.mockDb)

	ID := "random-id"
	handler, err := NewCmdMessageHandler(ID, suite.mockMessenger, suite.mockHub, suite.DB)
	require.Nil(suite.T(), err)

	msg := viewmodels.MessageView{
		Text: "This is a common message and should be ignored",
	}
	err = handler.HandleMessage(msg)
	assert.NoError(suite.T(), err)

	assert.NoError(suite.T(), suite.mockDb.ExpectationsWereMet())
	suite.mockHub.AssertExpectations(suite.T())
	suite.mockMessenger.AssertExpectations(suite.T())
}

func (suite *CommandMessageHandlerSuite) TestHandleMessageCommand() {
	expectTestUser(suite.mockDb)
	suite.mockMessenger.On("Publish", "100", "AAPL").Once()

	ID := "random-id"
	handler, err := NewCmdMessageHandler(ID, suite.mockMessenger, suite.mockHub, suite.DB)
	require.Nil(suite.T(), err)

	msg := viewmodels.MessageView{
		RoomID: 100,
		Text:   "/stock=AAPL",
	}
	err = handler.HandleMessage(msg)
	assert.NoError(suite.T(), err)

	assert.NoError(suite.T(), suite.mockDb.ExpectationsWereMet())
	suite.mockHub.AssertExpectations(suite.T())
	suite.mockMessenger.AssertExpectations(suite.T())
}

func (suite *CommandMessageHandlerSuite) TestCmdResponseHandler() {
	expectTestUser(suite.mockDb)
	suite.mockDb.ExpectBegin()
	suite.mockDb.ExpectQuery(`INSERT INTO "messages" (.+)`).
		WillReturnRows(sqlmock.NewRows([]string{"a"}).AddRow(1).AddRow(1))
	suite.mockDb.ExpectCommit()

	suite.mockHub.On("BroadcastMessage", mock.AnythingOfType("viewmodels.MessageView")).
		Return().Once()

	ID := "random-id"
	handler, err := NewCmdMessageHandler(ID, suite.mockMessenger, suite.mockHub, suite.DB)
	require.Nil(suite.T(), err)

	msg := messenger.BotMessage{
		RoomID:  10,
		Message: "Bot Message",
	}
	err = handler.CmdResponseHandler(msg)
	assert.NoError(suite.T(), err)

	assert.NoError(suite.T(), suite.mockDb.ExpectationsWereMet())
	suite.mockHub.AssertExpectations(suite.T())
	suite.mockMessenger.AssertExpectations(suite.T())
}

var userRows = []string{"username", "password", "email", "first_name", "last_name"}

func expectTestUser(mockDb sqlmock.Sqlmock) {
	mockDb.ExpectQuery(`SELECT (.+) FROM "users" (.+)`).
		WillReturnRows(sqlmock.NewRows(userRows).
			AddRow("Bot", "password", "bot@email.com", "First", "Last"))
}

type MockBotCommandMessenger struct {
	mock.Mock
}

func NewMockBotCommandMessenger() *MockBotCommandMessenger {
	return &MockBotCommandMessenger{}
}

func (m *MockBotCommandMessenger) Publish(key, message string) error {
	m.Called(key, message)
	return nil
}

func (m *MockBotCommandMessenger) StartHandler(func(string) (string, error)) error {
	return nil
}

func (m *MockBotCommandMessenger) StartConsumer(func(messenger.BotMessage) error) error {
	return nil
}
