package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrations(t *testing.T) {
	withRedis(t, "migration", func(t *testing.T, store Store) {
		bg := context.Background()

		t.Run("Migrations", func(t *testing.T) {
			_ = store.Flush(bg)
			_, err := store.GetQueue(bg, "default")
			assert.NoError(t, err)

			r := store.Redis()
			// create a queue with the old naming
			val := r.LPush(bg, "default", "mike").Val()
			assert.EqualValues(t, 1, val)

			assert.EqualValues(t, 1, r.Exists(bg, "default").Val())
			assert.EqualValues(t, 0, r.Exists(bg, "q:default").Val())

			ver, err := store.DataVersion(bg)
			assert.NoError(t, err)
			assert.EqualValues(t, 0, ver)

			ver, err = store.ApplyMigrations(bg)
			assert.NoError(t, err)
			assert.EqualValues(t, 1, ver)

			assert.EqualValues(t, 0, r.Exists(bg, "default").Val())
			assert.EqualValues(t, 1, r.Exists(bg, "q:default").Val())
		})
	})
}
