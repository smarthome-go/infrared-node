## Changelog for v1.4.0

### Token Authentication
- Added support for token authentication (*which is supported since server version `0.0.55`*)
- Improved several server logs

#### Example Configuration File
- Once you set `tokenAuth` to `true`, any specified username or password will be ignored (*token has precedence*)
- A token should identify a user who has at least the permission *homescript*
- A token can be generated under `http://smarthome.box/profile` (*substitute `http://smarthome.box` with your domain or IP*)

```json
{
	"smarthome": {
		"url": "http://smarthome.box",
		"tokenAuth": true,
		"credentials": {
			"user": "",
			"password": "",
			"token": "your-token-here"
		},
		"hmsTimeout": 10
	},
	"hardware": {
		"enabled": false,
		"pin": 0
	},
	"actions": [
		{
			"name": "demo",
			"code": "2a00aaa95",
			"homescript": "switch('sx', on)"
		}
	]
}
```


