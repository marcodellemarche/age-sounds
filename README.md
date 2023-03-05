# Age of Empires 2 sounds

This project is meant to download all the Age of Empires 2 sounds, from the Famdom website.

It is written in Go and uses Colly as the web scraper.

## Convert OGG files to WAV files

For the time being, the sound files are downloaded in `.ogg` format. However, the [Random Notification Sounds](https://play.google.com/store/apps/details?id=com.simplycomplexapps.randomnotificationsounds) Android app which will digest these files, wants them in `.wav` format. The following bash script can do the conversion:

```bash
cd download
mkdir wav

for i in *.ogg; do     
  ffmpeg -acodec libvorbis -i "$i" -acodec pcm_s16le "wav/${i%ogg}wav"
done
```
