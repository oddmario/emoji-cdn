package utils

import (
	"errors"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/tidwall/gjson"
)

func EmojipediaScraper(emoji, style string) (io.ReadCloser, string, error) {
	var emojipedia_nextjs_key string = "Ac4U6X8McKS9ckPowL50Y"
	var emojipedia_cdn_base string = "https://em-content.zobj.net/"

	var get_emoji_slug_body string = ""

	for range 5 {
		get_emoji_slug, err := HttpClient.R().
			SetHeaders(map[string]string{
				"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0",
				"Accept":          "*/*",
				"Accept-Language": "en-US,en;q=0.5",
				"Referer":         "https://emojipedia.org/",
				"x-nextjs-data":   "1",
				"Connection":      "keep-alive",
				"Sec-Fetch-Dest":  "empty",
				"Sec-Fetch-Mode":  "cors",
				"Sec-Fetch-Site":  "same-origin",
				"Priority":        "u=0",
				"Pragma":          "no-cache",
				"Cache-Control":   "no-cache",
				"TE":              "trailers",
			}).
			Get("https://emojipedia.org/_next/data/" + emojipedia_nextjs_key + "/en/search.json?q=" + url.QueryEscape(emoji))

		if err != nil {
			return nil, "", errors.New("an internal error has occurred")
		}

		get_emoji_slug_body = get_emoji_slug.String()

		if !strings.Contains(get_emoji_slug_body, "__N_REDIRECT") {
			continue
		} else {
			break
		}
	}

	if !strings.Contains(get_emoji_slug_body, "__N_REDIRECT") {
		return nil, "", errors.New("invalid emojipedia response")
	}

	get_emoji_slug_parser := gjson.Parse(get_emoji_slug_body)

	emoji_slug := strings.Split(get_emoji_slug_parser.Get("pageProps").Get("__N_REDIRECT").String(), "/")[1]

	var get_emoji_body string = ""

	for range 5 {
		get_emoji, err := HttpClient.R().
			SetHeaders(map[string]string{
				"User-Agent":      "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0",
				"Accept":          "*/*",
				"Accept-Language": "en-US,en;q=0.5",
				"Referer":         "https://emojipedia.org/" + emoji_slug,
				"x-nextjs-data":   "1",
				"Connection":      "keep-alive",
				"Sec-Fetch-Dest":  "empty",
				"Sec-Fetch-Mode":  "cors",
				"Sec-Fetch-Site":  "same-origin",
				"Priority":        "u=1, i",
				"Pragma":          "no-cache",
				"Cache-Control":   "no-cache",
				"TE":              "trailers",
				"x-client":        "emojipedia.org",
				"x-query-hash":    "237744ac6bfd1775eca976e362f4996cf0dd41cd",
				"Content-Type":    "application/json",
			}).
			SetBody(map[string]interface{}{
				"operationName": "emojiV1",
				"query":         "\n    query emojiV1($slug: Slug!, $lang: Language) {\n      emoji_v1(slug: $slug, lang: $lang) {\n        ...emojiDetailsResource\n      }\n    }\n    \n  fragment vendorAndPlatformResource on VendorAndPlatform {\n    slug\n    title\n    description\n    items {\n      date\n      slug\n      title\n      image {\n        source\n        description\n        useOriginalImage\n      }\n    }\n  }\n\n    \n  fragment shortCodeResource on Shortcode {\n    code\n    vendor {\n      slug\n      title\n    }\n  }\n\n    \n  fragment emojiResource on Emoji {\n    id\n    title\n    code\n    slug\n    currentCldrName\n    codepointsHex\n    description\n    modifiers\n    appleName\n    alsoKnownAs\n    shortcodes {\n      ...shortCodeResource\n    }\n    proposals {\n      name\n      url\n    }\n  }\n\n    \n  fragment emojiDetailsResource on Emoji {\n    ...emojiResource\n    type\n    socialImage {\n      source\n    }\n    emojiVersion {\n      ...emojiVersionResource\n    }\n    version {\n      ...versionUnicodeResource\n    }\n    components {\n      ...emojiResource\n    }\n    goesGreatWith {\n      ... on Emoji {\n        title\n        slug\n        code\n        currentCldrName\n        description\n      }\n      ... on StaticContent {\n        title\n        slug\n        titleEmoji {\n          code\n          title\n          currentCldrName\n          description\n          slug\n        }\n      }\n    }\n    genderVariants {\n      slug\n      title\n      currentCldrName\n    }\n    skinTone\n    skinToneVariants {\n      slug\n      skinTone\n      code\n      title\n      currentCldrName\n    }\n    vendorsAndPlatforms {\n      ...vendorAndPlatformResource\n    }\n    alerts {\n      name\n      content\n      notes\n      link\n      representingEmoji {\n        code\n      }\n    }\n    alert {\n      content\n      link\n      representingEmoji {\n        code\n        slug\n      }\n    }\n    shortcodes {\n      code\n      source\n      vendor {\n        slug\n        title\n      }\n    }\n    shopInfo {\n      ... on ShopInfo {\n        url\n        image\n      }\n    }\n    mashup {\n      supported\n    }\n  }\n\n    \n  fragment emojiVersionResource on EmojiVersion {\n    name\n    date\n    slug\n    status\n  }\n\n    \n  fragment versionUnicodeResource on Unicode {\n    name\n    slug\n    date\n    description\n    status\n  }\n\n  ",
				"variables": map[string]string{
					"lang": "EN",
					"slug": emoji_slug,
				},
			}).
			Post("https://emojipedia.org/api/graphql")

		if err != nil {
			return nil, "", errors.New("an internal error has occurred")
		}

		get_emoji_body = get_emoji.String()

		if !strings.Contains(get_emoji_body, emoji_slug) {
			continue
		} else {
			break
		}
	}

	if !strings.Contains(get_emoji_body, emoji_slug) {
		return nil, "", errors.New("invalid emojipedia response")
	}

	emoji_data := gjson.Parse(get_emoji_body).Get("data").Get("emoji_v1")

	platformIndex := int64(-1)

	emoji_data.Get("vendorsAndPlatforms").ForEach(func(key, value gjson.Result) bool {
		if value.Get("slug").String() == style {
			platformIndex = key.Int()

			return false
		}

		return true // keep iterating
	})

	if platformIndex < 0 {
		return nil, "", errors.New("no emoji found for the specified style")
	}

	emoji_img := emojipedia_cdn_base + emoji_data.Get("vendorsAndPlatforms").Get(I64ToStr(platformIndex)).Get("items").Get("0").Get("image").Get("source").String()

	emoji_proxy, err := HttpClient.R().
		SetDoNotParseResponse(true).
		Get(emoji_img)

	if err != nil {
		return nil, "", errors.New("an internal error has occurred")
	}

	u, _ := url.Parse(emoji_img)
	ext := path.Ext(u.Path)

	reader := emoji_proxy.RawResponse.Body

	return reader, ext, nil
}
