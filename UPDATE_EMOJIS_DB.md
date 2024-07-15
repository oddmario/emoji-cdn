# 🧐 Update/create the Emojis database

## Stage I. Preparing the `emojis.txt` file

In order to update or create the emojis database, this application requires an `emojis.txt` file containing a list of the emojis that it should scrape from Emojipedia.

The `emojis.txt` file should be placed beside the binary file of the **emoji-cdn** application, and it should contain the emojis one per line.

For example:
```txt
🫣
😞
```

Below is one way to create an emojis.txt file containing all the emojis available:
1. Go to https://getemoji.com/
2. Perform any user interacion activity on the page (e.g. click anywhere on the page for example)
3. Run the following snippet in the console of your browser:
```js
var finalOutput = "";

document.querySelectorAll(".emoji-large").forEach((emoji) => {
    finalOutput += emoji.textContent + "\n";
});

navigator.clipboard.writeText(finalOutput);
```
4. Voila! You have the emojis list on your clipboard now.
5. Paste the clipboard into a file named `emojis.txt`
6. Make sure the `emojis.txt` file is placed beside the binary file of **emoji-cdn**.

## Stage II. Starting the scraper
To start scraping from Emojipedia, run **emoji-cdn** in the scraping mode:
```
./emoji-cdn --update-db
```

To speed things up, the scraper scrapes in multiple threads. The default is 10 threads.

To scrape faster, try a higher number of threads like 50:
```
./emoji-cdn --update-db --threads=50
```

Kindly note that the scraper may take some time even with a high number of threads. This happens due to the nature of the Emojipedia server response times (they are sometimes not the fastest), and also because some emoji images exceed the size of 500 KB.

If you prefer starting the scraper in the background (for example, inside a `screen` session), feel free to do that :)

### And tada!
That's it really!

Once the scraper finishes its job, you shall find all the scraped emojis in a directory named `.emojis-db` placed alongside the **emoji-cdn** binary.

**emoji-cdn** will use the data inside that directory to serve its requests.