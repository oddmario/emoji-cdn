# ⚡️ Emoji CDN
## a blazing fast emojis CDN that supports 30 emoji sets

![GitHub License](https://img.shields.io/github/license/oddmario/emoji-cdn)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/oddmario/emoji-cdn)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues/oddmario/emoji-cdn)
![GitHub Issues or Pull Requests](https://img.shields.io/github/issues-pr/oddmario/emoji-cdn)

<h1 align="center">
  🤩 🕺 🥳 
</h1>

-----

## Usage
Simply, use the base URL for **emoji-cdn** which is `https://emoji-cdn.mqrio.dev/[emoji]?style=[style]`

The CDN will respond with an image containing the emoji for the specified style (emoji set/design)

For example:
```html
<img src="https://emoji-cdn.mqrio.dev/🎉?style=google" />
```

## Supported styles
- apple
- google
- samsung
- microsoft-3D-fluent
- microsoft
- whatsapp
- twitter
- facebook
- huawei
- joypixels
- lg
- telegram
- animated-noto-color-emoji
- microsoft-teams
- skype
- twitter-emoji-stickers
- joypixels-animations
- serenityos
- toss-face
- sony
- noto-emoji
- openmoji
- icons8
- emojidex
- messenger
- htc
- softbank
- docomo
- au-kddi
- mozilla

## Why not [benborgers/emojicdn](https://github.com/benborgers/emojicdn)?
The emojicdn project by Ben Borgers is definitely a cool project! It actually is what inspired me to create this whole project.

However, I found myself in the need of a little faster serving speed (despite how the benborgers CDN is fast enough).

This project stores and fetches the emojis locally, unlike `benborgers/emojicdn` which behaves as a reverse proxy for a remotely hosted emojis source.

In addition to that, I also needed to use some other styles which were missing from the benborgers project (especially the WhatsApp emoji set).

So in order to make up for all that, this project was born.

## Self Hosting
You are absolutely more than welcome to self-host your own instance of **emoji-cdn** :)

Here is how:

1. Build the project:
```
python3 build.py
```

2. [Create your own Emojis database](/UPDATE_EMOJIS_DB.md)

3. Run the application to start emoting!

## License
- Code licensed under the MIT License: http://opensource.org/licenses/MIT
- All the emojis were obtained from [Emojipedia.org](https://emojipedia.org/) and are owned by their original distributors & designers