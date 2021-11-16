package datastore_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/asdine/storm/v3/q"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/ogma/cmd"
	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

func TestManagerNew(t *testing.T) {
	manager, dbFilePath, err := initDatastoreManager()
	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()
	defer manager.Stop()

	assert.NoError(t, err)

	_, err = os.Stat(dbFilePath)

	assert.NotErrorIs(t, err, os.ErrNotExist)
}

func TestManagerNewFail(t *testing.T) {
	manager, dbFilePath, err := initDatastoreManager()
	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()
	defer manager.Stop()

	_, err = datastore.New(dbFilePath)
	assert.Error(t, err, "datastore manager should fail to open duplicate db file")
}

func TestManagerGetPath(t *testing.T) {
	manager, dbFilePath, err := initDatastoreManager()
	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()
	defer manager.Stop()

	assert.Equal(t, dbFilePath, manager.GetPath())
}

func initDatastoreManager() (*datastore.Manager, string, error) {
	currentTime := time.Now()
	filename := fmt.Sprintf("test_%d.db", currentTime.Unix())
	manager, err := datastore.New(filename)

	return manager, filename, err
}

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestManager_Save(t *testing.T) {
	manager, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer manager.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	type testData struct {
		ID      int
		TestNum int
	}

	td := testData{ID: 1, TestNum: 5}
	err = manager.Save(&td)
	assert.NoError(t, err)
}

func TestManager_One(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer m.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Mollit",
			PageNumber:          1,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 1234,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Eiusmod",
			PageNumber:          2,
			IndexedCategory:     "Commodo",
			IndexedMemberNumber: 1234,
			MemberExtension:     "B",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Magna officia anim dolore enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	type args struct {
		fieldName string
		value     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty search",
			args: args{
				fieldName: "",
				value:     nil,
			},
			wantErr: true,
		},
		{
			name: "good - with results",
			args: args{
				fieldName: "IndexedMemberNumber",
				value:     1234,
			},
			wantErr: false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				fieldName: "IndexedMemberNumber",
				value:     5678,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got lstg.Listing

			if err := m.One(tt.args.fieldName, tt.args.value, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.One() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Find(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer m.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Mollit",
			PageNumber:          1,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 1234,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Eiusmod",
			PageNumber:          2,
			IndexedCategory:     "Commodo",
			IndexedMemberNumber: 1234,
			MemberExtension:     "B",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Magna officia anim dolore enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Id",
			PageNumber:          3,
			IndexedCategory:     "Consequat",
			IndexedMemberNumber: 5678,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Velit cillum cillum ea officia nulla enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	type args struct {
		fieldName string
		value     interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "empty search",
			args: args{
				fieldName: "",
				value:     nil,
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "good - with results",
			args: args{
				fieldName: "IndexedMemberNumber",
				value:     1234,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				fieldName: "IndexedMemberNumber",
				value:     1,
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []lstg.Listing

			if err := m.Find(tt.args.fieldName, tt.args.value, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Find() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_AllByIndex(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer m.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Mollit",
			PageNumber:          1,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 1234,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Eiusmod",
			PageNumber:          2,
			IndexedCategory:     "Commodo",
			IndexedMemberNumber: 1234,
			MemberExtension:     "B",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Magna officia anim dolore enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Id",
			PageNumber:          3,
			IndexedCategory:     "Consequat",
			IndexedMemberNumber: 5678,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Velit cillum cillum ea officia nulla enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	type args struct {
		fieldName string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{

		{
			name: "good - with results",
			args: args{
				fieldName: "IndexedMemberNumber",
			},
			wantCount: 3,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []lstg.Listing

			if err := m.AllByIndex(tt.args.fieldName, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.AllByIndex() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_All(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer m.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Mollit",
			PageNumber:          1,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 1234,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Eiusmod",
			PageNumber:          2,
			IndexedCategory:     "Commodo",
			IndexedMemberNumber: 1234,
			MemberExtension:     "B",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Magna officia anim dolore enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Id",
			PageNumber:          3,
			IndexedCategory:     "Consequat",
			IndexedMemberNumber: 5678,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Velit cillum cillum ea officia nulla enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{
			name: "good - with results",

			wantCount: 3,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []lstg.Listing

			if err := m.All(&got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.All() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_Select(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer m.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Mollit",
			PageNumber:          1,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 1234,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Eiusmod",
			PageNumber:          2,
			IndexedCategory:     "Commodo",
			IndexedMemberNumber: 1234,
			MemberExtension:     "B",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Magna officia anim dolore enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Id",
			PageNumber:          3,
			IndexedCategory:     "Consequat",
			IndexedMemberNumber: 5678,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Velit cillum cillum ea officia nulla enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	type args struct {
		matcher q.Matcher
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "good - with results",
			args: args{
				matcher: q.Eq("IndexedMemberNumber", 1234),
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				matcher: q.Eq("IndexedMemberNumber", 1),
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []lstg.Listing

			q := m.Select(tt.args.matcher)
			if err := q.Find(&got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Select() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_Range(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer m.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Mollit",
			PageNumber:          1,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 1234,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Eiusmod",
			PageNumber:          2,
			IndexedCategory:     "Commodo",
			IndexedMemberNumber: 1234,
			MemberExtension:     "B",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Magna officia anim dolore enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Id",
			PageNumber:          3,
			IndexedCategory:     "Consequat",
			IndexedMemberNumber: 5678,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Velit cillum cillum ea officia nulla enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	type args struct {
		fieldName string
		min       interface{}
		max       interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "empty search",
			args: args{
				fieldName: "",
				min:       nil,
				max:       nil,
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "good - subset results",
			args: args{
				fieldName: "IndexedMemberNumber",
				min:       1235,
				max:       5679,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "good - full results", // check for specific error "not found"
			args: args{
				fieldName: "IndexedMemberNumber",
				min:       1,
				max:       9999,
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				fieldName: "IndexedMemberNumber",
				min:       nil,
				max:       nil,
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []lstg.Listing

			if err := m.Range(tt.args.fieldName, tt.args.min, tt.args.max, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Range() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_Prefix(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer m.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Mollit",
			PageNumber:          1,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 1234,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Eiusmod",
			PageNumber:          2,
			IndexedCategory:     "Commodo",
			IndexedMemberNumber: 1234,
			MemberExtension:     "B",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Magna officia anim dolore enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Id",
			PageNumber:          3,
			IndexedCategory:     "Consequat",
			IndexedMemberNumber: 5678,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Velit cillum cillum ea officia nulla enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	type args struct {
		fieldName string
		prefix    string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "good - with results",
			args: args{
				fieldName: "IndexedCategory",
				prefix:    "Co",
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				fieldName: "IndexedCategory",
				prefix:    "Zzz",
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []lstg.Listing

			if err := m.Prefix(tt.args.fieldName, tt.args.prefix, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Prefix() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_Count(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer m.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Mollit",
			PageNumber:          1,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 1234,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Eiusmod",
			PageNumber:          2,
			IndexedCategory:     "Commodo",
			IndexedMemberNumber: 1234,
			MemberExtension:     "B",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Magna officia anim dolore enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Id",
			PageNumber:          3,
			IndexedCategory:     "Consequat",
			IndexedMemberNumber: 5678,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Velit cillum cillum ea officia nulla enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "count all listings",
			wantCount: 3,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got lstg.Listing

			c, err := m.Count(&got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.Find() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, c, "Found %d results, wanted %d", c, tt.wantCount)
		})
	}
}
