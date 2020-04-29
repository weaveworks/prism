package local

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cortexproject/cortex/pkg/chunk/util"
)

func TestFSObjectClient_DeleteChunksBefore(t *testing.T) {
	deleteFilesOlderThan := 10 * time.Minute

	fsChunksDir, err := ioutil.TempDir(os.TempDir(), "fs-chunks")
	require.NoError(t, err)

	bucketClient, err := NewFSObjectClient(FSConfig{
		Directory: fsChunksDir,
	})
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(fsChunksDir))
	}()

	file1 := "file1"
	file2 := "file2"

	// Creating dummy files
	require.NoError(t, os.Chdir(fsChunksDir))

	f, err := os.Create(file1)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	f, err = os.Create(file2)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	// Verify whether all files are created
	files, _ := ioutil.ReadDir(".")
	require.Equal(t, 2, len(files), "Number of files should be 2")

	// No files should be deleted, since all of them are not much older
	require.NoError(t, bucketClient.DeleteChunksBefore(context.Background(), time.Now().Add(-deleteFilesOlderThan)))
	files, _ = ioutil.ReadDir(".")
	require.Equal(t, 2, len(files), "Number of files should be 2")

	// Changing mtime of file1 to make it look older
	require.NoError(t, os.Chtimes(file1, time.Now().Add(-deleteFilesOlderThan), time.Now().Add(-deleteFilesOlderThan)))
	require.NoError(t, bucketClient.DeleteChunksBefore(context.Background(), time.Now().Add(-deleteFilesOlderThan)))

	// Verifying whether older file got deleted
	files, _ = ioutil.ReadDir(".")
	require.Equal(t, 1, len(files), "Number of files should be 1 after enforcing retention")
}

func TestFSObjectClient_List(t *testing.T) {
	fsObjectsDir, err := ioutil.TempDir(os.TempDir(), "fs-objects")
	require.NoError(t, err)

	bucketClient, err := NewFSObjectClient(FSConfig{
		Directory: fsObjectsDir,
	})
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(fsObjectsDir))
	}()

	foldersWithFiles := make(map[string][]string)
	foldersWithFiles["folder1/"] = []string{"file1", "file2"}
	foldersWithFiles["folder2/"] = []string{"file3", "file4", "file5"}

	for folder, files := range foldersWithFiles {
		for _, filename := range files {
			err := bucketClient.PutObject(context.Background(), folder+filename, bytes.NewReader([]byte(filename)))
			require.NoError(t, err)
		}
	}

	// create an empty directory which should get excluded from the list
	require.NoError(t, util.EnsureDirectory(filepath.Join(fsObjectsDir, "empty-folder")))

	files := []string{"outer-file1", "outer-file2"}

	for _, fl := range files {
		err := bucketClient.PutObject(context.Background(), fl, bytes.NewReader([]byte(fl)))
		require.NoError(t, err)
	}

	storageObjects, commonPrefixes, err := bucketClient.List(context.Background(), "")
	require.NoError(t, err)

	require.Len(t, storageObjects, len(files))
	for i := range storageObjects {
		require.Equal(t, storageObjects[i].Key, files[i])
	}

	require.Len(t, commonPrefixes, len(foldersWithFiles))
	for _, commonPrefix := range commonPrefixes {
		_, ok := foldersWithFiles[string(commonPrefix)]
		require.True(t, ok)
	}

	for folder, files := range foldersWithFiles {
		storageObjects, commonPrefixes, err := bucketClient.List(context.Background(), folder)
		require.NoError(t, err)

		require.Len(t, storageObjects, len(files))
		for i := range storageObjects {
			require.Equal(t, storageObjects[i].Key, folder+files[i])
		}

		require.Len(t, commonPrefixes, 0)
	}
}

func TestFSObjectClient_DeleteObject(t *testing.T) {
	fsObjectsDir, err := ioutil.TempDir(os.TempDir(), "fs-delete-object")
	require.NoError(t, err)

	bucketClient, err := NewFSObjectClient(FSConfig{
		Directory: fsObjectsDir,
	})
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(fsObjectsDir))
	}()

	foldersWithFiles := make(map[string][]string)
	foldersWithFiles["folder1/"] = []string{"file1", "file2"}

	for folder, files := range foldersWithFiles {
		for _, filename := range files {
			err := bucketClient.PutObject(context.Background(), folder+filename, bytes.NewReader([]byte(filename)))
			require.NoError(t, err)
		}
	}

	// let us check if we have right folders created
	_, commonPrefixes, err := bucketClient.List(context.Background(), "")
	require.NoError(t, err)
	require.Len(t, commonPrefixes, len(foldersWithFiles))

	// let us delete file1 from folder1 and check that file1 is gone but folder1 with file2 is still there
	require.NoError(t, bucketClient.DeleteObject(context.Background(), "folder1/file1"))
	_, err = os.Stat(filepath.Join(fsObjectsDir, "folder1/file1"))
	require.True(t, os.IsNotExist(err))

	_, err = os.Stat(filepath.Join(fsObjectsDir, "folder1/file2"))
	require.NoError(t, err)

	// let us delete second file as well and check that folder1 also got removed
	require.NoError(t, bucketClient.DeleteObject(context.Background(), "folder1/file2"))
	_, err = os.Stat(filepath.Join(fsObjectsDir, "folder1"))
	require.True(t, os.IsNotExist(err))

	_, err = os.Stat(fsObjectsDir)
	require.NoError(t, err)

	// let us see ensure folder2 is still there will all the files:
	/*files, commonPrefixes, err := bucketClient.List(context.Background(), "folder2/")
	require.NoError(t, err)
	require.Len(t, commonPrefixes, 0)
	require.Len(t, files, len(foldersWithFiles["folder2/"]))*/
}

func TestIsNotEmptyErr(t *testing.T) {
	outerDir, err := ioutil.TempDir(os.TempDir(), "empty-dir-err-check")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(outerDir))
	}()

	// create a directory inside the outer dir
	require.NoError(t, os.Mkdir(filepath.Join(outerDir, "nested"), 0700))

	// try removing the outer directory and see if it throws syscall.ENOTEMPTY error
	err = os.Remove(outerDir)
	require.Error(t, err)
	require.True(t, isNotEmptyErr(err))

	// try removing a non-existent directory and see if it does not throw syscall.ENOTEMPTY error
	err = os.Remove(filepath.Join(outerDir, "non-existent"))
	require.Error(t, err)
	require.False(t, isNotEmptyErr(err))
}
