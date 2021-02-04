# Secret Music Exchange!

Dear {{.GiverName}},

You will be selecting music for **{{.ReceiverName}}**.

Here's what they had to say about themselves:

<dl>
{{ range .ReceiverResponses }}
  <dt>{{ .Question }}</dt>
  <dd>{{ .Answer }}</dd><br/>
{{ end }}
</dl>

### Instructions

##### 1: Select Music
Create a sharable playlist (between 45 and 75 minutes) on one of the platforms your recipient uses:

{{ range .ReceiverPlatforms }}- {{ . }}
{{ end }}
##### 2: Copy Link
If you are using an album, copy the album link from Spotify. If you are creating a playlist, copy the playlist link.

##### 3: Share Music
We'll do this all at once. When it's time, you will share that link on Slack. @mention then in the #music-appreciation channel, so we can all see what you picked.

##### 4: Enjoy the music shared with you!
You'll be getting some music to enjoy too.
