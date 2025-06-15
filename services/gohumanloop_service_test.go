package services

import (
	"testing"
	"time"

	"github.com/ptonlix/gohumanloop-wework/models"
	"github.com/stretchr/testify/assert"
)

func TestGoHumanloopService(t *testing.T) {
	service := NewGoHumanloopService()

	// Setup test data
	setupTestData := func() {
		// Clear existing data
		service.ormer.Raw("DELETE FROM human_loop").Exec()

		// Create test records
		testRecords := []*models.HumanLoop{
			{
				TaskId:         "task-1",
				ConversationId: "conv-1",
				RequestId:      "req-1",
				LoopType:       "type-1",
				Platform:       "test",
				Status:         "pending",
				Created:        time.Now(),
			},
			{
				TaskId:         "task-2",
				ConversationId: "conv-1",
				RequestId:      "req-2",
				LoopType:       "type-2",
				Platform:       "test",
				Status:         "pending",
				Created:        time.Now().Add(time.Minute),
			},
			{
				TaskId:         "task-3",
				ConversationId: "conv-2",
				RequestId:      "req-3",
				LoopType:       "type-3",
				Platform:       "test",
				Status:         "pending",
				Created:        time.Now().Add(2 * time.Minute),
			},
		}

		for _, record := range testRecords {
			service.CreateHumanLoop(record)
		}
	}

	t.Run("TestCreateHumanLoop", func(t *testing.T) {
		hl := &models.HumanLoop{
			TaskId:         "test-task",
			ConversationId: "test-conv",
			RequestId:      "test-req",
			LoopType:       "test",
			Platform:       "test",
			Status:         "pending",
		}

		id, err := service.CreateHumanLoop(hl)
		assert.NoError(t, err)
		assert.True(t, id > 0)
	})

	t.Run("TestGetHumanLoopByID", func(t *testing.T) {
		setupTestData()
		hl, err := service.GetHumanLoopByID(1)
		assert.NoError(t, err)
		assert.Equal(t, "test-task", hl.TaskId)
	})

	t.Run("TestGetHumanLoopByRequestID", func(t *testing.T) {
		setupTestData()
		hl, err := service.GetHumanLoopByRequestId("conv-1", "req-1", "test")
		assert.NoError(t, err)
		assert.Equal(t, "task-1", hl.TaskId)
	})

	t.Run("TestUpdateHumanLoop", func(t *testing.T) {
		setupTestData()
		hl := &models.HumanLoop{
			ID:     1,
			Status: "completed",
		}

		rows, err := service.UpdateHumanLoop(hl)
		assert.NoError(t, err)
		assert.True(t, rows > 0)

		// Verify update
		updatedHl, _ := service.GetHumanLoopByID(1)
		assert.Equal(t, "completed", updatedHl.Status)
	})

	t.Run("TestListHumanLoops", func(t *testing.T) {
		setupTestData()
		filter := map[string]interface{}{"status": "pending"}
		hls, total, err := service.ListHumanLoops(filter, 1, 10)
		assert.NoError(t, err)
		assert.True(t, total >= 3)
		assert.Equal(t, "pending", hls[0].Status)
	})

	t.Run("TestDeleteHumanLoop", func(t *testing.T) {
		setupTestData()
		err := service.DeleteHumanLoop(1)
		assert.NoError(t, err)

		// Verify soft delete
		hl, err := service.GetHumanLoopByID(1)
		assert.NoError(t, err)
		assert.True(t, hl.IsDeleted)
	})

	t.Run("TestGetHumanLoopsByConversationId", func(t *testing.T) {
		setupTestData()
		// Test getting records for conv-1
		hls, err := service.GetHumanLoopsByConversationId("conv-1", "test")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(hls))
		assert.Equal(t, "conv-1", hls[0].ConversationId)
		assert.Equal(t, "conv-1", hls[1].ConversationId)

		// Verify ordering (newest first)
		assert.True(t, hls[0].Created.After(hls[1].Created))

		// Test getting records for non-existent conversation
		hls, err = service.GetHumanLoopsByConversationId("nonexistent", "test")
		assert.NoError(t, err)
		assert.Equal(t, 0, len(hls))
	})

	t.Run("TestBatchUpdateHumanLoops", func(t *testing.T) {
		setupTestData()

		// Get existing records to update
		hls, _ := service.GetHumanLoopsByConversationId("conv-1", "test")
		assert.Equal(t, 2, len(hls))

		// Modify records
		for i := range hls {
			hls[i].Status = "completed"
			hls[i].TaskId = "updated-task"
		}

		// Perform batch update
		err := service.BatchUpdateHumanLoops(hls)
		assert.NoError(t, err)

		// Verify updates
		updatedHls, _ := service.GetHumanLoopsByConversationId("conv-1", "test")
		for _, hl := range updatedHls {
			assert.Equal(t, "completed", hl.Status)
			assert.Equal(t, "updated-task", hl.TaskId)
		}

		// Test transaction rollback on error
		invalidHls := []*models.HumanLoop{
			{ID: 9999, Status: "should-fail"}, // Non-existent ID
		}
		err = service.BatchUpdateHumanLoops(invalidHls)
		assert.Error(t, err)

		// Verify no records were updated by the failed transaction
		unchangedHl, _ := service.GetHumanLoopByID(1)
		assert.Equal(t, "completed", unchangedHl.Status) // Still has previous update
	})
}
