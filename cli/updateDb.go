package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"mario/emoji-cdn/constants"
	"mario/emoji-cdn/lock"
	"mario/emoji-cdn/utils"
	"os"
	"path/filepath"
	"time"
)

func updateDbThread(emoji string, active_threads_count_lock, emojis_count_lock *lock.ThreadsafeVariable) {
	active_threads_count_lock.Lock()
	active_threads_count_lock.Set(int(active_threads_count_lock.Get().(int) + 1))
	active_threads_count_lock.Unlock()

	fmt.Println("started scraper for", emoji)

	emojisDbPath := filepath.Join(".", ".emojis-db", emoji)
	err := os.MkdirAll(emojisDbPath, os.ModePerm)
	if err != nil {
		active_threads_count_lock.Lock()
		active_threads_count_lock.Set(int(active_threads_count_lock.Get().(int) - 1))
		active_threads_count_lock.Unlock()

		emojis_count_lock.Lock()
		emojis_count_lock.Set(int(emojis_count_lock.Get().(int) + 1))
		emojis_count_lock.Unlock()

		return
	}

	readers, err := utils.EmojipediaScraper(emoji, constants.EmojipediaPlatforms)

	if err != nil {
		fmt.Println("scrape err", emoji, err)

		active_threads_count_lock.Lock()
		active_threads_count_lock.Set(int(active_threads_count_lock.Get().(int) - 1))
		active_threads_count_lock.Unlock()

		emojis_count_lock.Lock()
		emojis_count_lock.Set(int(emojis_count_lock.Get().(int) + 1))
		emojis_count_lock.Unlock()

		return
	}

	for filename, reader := range readers {
		f, err := os.OpenFile(filepath.Join(".", ".emojis-db", emoji, filename), os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			fmt.Println("error while opening file", err)

			reader.Close()

			continue
		}

		io.Copy(f, reader)

		f.Close()
		reader.Close()

		fmt.Println(emoji, filename, "has been scraped")
	}

	active_threads_count_lock.Lock()
	active_threads_count_lock.Set(int(active_threads_count_lock.Get().(int) - 1))
	active_threads_count_lock.Unlock()

	emojis_count_lock.Lock()
	emojis_count_lock.Set(int(emojis_count_lock.Get().(int) + 1))
	emojis_count_lock.Unlock()
}

func UpdateDb(maxThreads int) {
	emojisFilename := "emojis.txt"
	path, _ := filepath.Abs("./" + emojisFilename)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Println("The `emojis.txt` file does not exist. Exiting...")

		return
	}

	emojisCount := 0

	fileEmojisCount, err := os.Open(path)
	if err != nil {
		log.Fatal(err)

		return
	}
	emojisCount, _ = utils.LineCounter(fileEmojisCount)
	fileEmojisCount.Close()

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)

		return
	}
	defer file.Close()

	initEmojisDb := filepath.Join(".", ".emojis-db")
	err = os.MkdirAll(initEmojisDb, os.ModePerm)
	if err != nil {
		log.Fatal(err)

		return
	}

	var activeThreadsCountLock *lock.ThreadsafeVariable = new(lock.ThreadsafeVariable)
	var emojisCountLock *lock.ThreadsafeVariable = new(lock.ThreadsafeVariable)

	activeThreadsCountLock.Set(int(0))
	emojisCountLock.Set(int(0))

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		for {
			activeThreadsCountLock.Lock()

			if activeThreadsCountLock.Get().(int) >= maxThreads {
				activeThreadsCountLock.Unlock()

				time.Sleep(1 * time.Second)

				continue
			} else {
				time.Sleep(5 * time.Millisecond)

				activeThreadsCountLock.Unlock()

				if activeThreadsCountLock.Get().(int) >= maxThreads {
					continue
				} else {
					break
				}
			}
		}

		go updateDbThread(scanner.Text(), activeThreadsCountLock, emojisCountLock)
	}

	for {
		if emojisCountLock.Get().(int) >= emojisCount {
			fmt.Println("Done.")

			break
		} else {
			time.Sleep(2 * time.Second)

			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
