# API Documentation

##### Table of Content

- [POST /auth](#post-auth)
- [POST /auth/logout](#post-authlogout)
- [GET /dataCollection/status](#get-datacollectionstatus)
- [GET /dataCollection/statistics](#get-datacollectionstatistics)
- [GET /sensors](#get-sensors)
- [GET /sensors/:sensor_id](#get-sensorssensor_id)
- [GET /sensors/:sensor_id/values](#get-sensorssensor_idvalues)
- [GET /sensors/:sensor_id/pushSettings [auth required]](#get-sensorssensor_idpushsettings-auth-required)
- [POST /sensors/:sensor_id/pushSettings [auth required]](#post-sensorssensor_idpushsettings-auth-required)
- [DELETE /sensors/:sensor_id/pushSettings/:id [auth required]](#delete-sensorssensor_idpushsettingsid-auth-required)
- [GET /myPushSettings/sensors [auth required]](#get-mypushsettingssensors-auth-required)
- [GET /search/sensors/:query](#get-searchsensorsquery)
- [GET /channels](#get-channels)
- [GET /channels/:channel_id](#get-channelschannel_id)
- [GET /channels/:channel_id/sensors](#get-channelschannel_idsensors)
- [GET /user](#get-user)
- [GET /userDevices](#get-userdevices)

### POST /auth

This API receives `username` and `password` of a [Waziup user](dashboard.waziup.io) and authenticates it through Waziup provided API and on success, it returns a token that can be used to call the protected APIs of the Simulator.

If you call this function through a browser, it automatically saves the token to a cookie, so you do not need to include it in every future calls until the cookie gets expired.

If no APIs are called in a specific period (_default is 20 minutes_), the cookie gets expired and the user needs to be authenticated again.

#### Input Format:

```
{
	"username": <String>,
	"password": <String>
}
```

#### Call Example:

```
curl -X POST -H 'Content-Type: application/json' -i http://localhost:8080/auth --data '{ "username": "cdupont", "password": "password"}'
```

**Output:**

```
$2a$10$4LYXlpTVwUM73DFrtPEWZOE2IEglleW75.3OYO2OkfQ90jN7iI90S
```

---

### POST /auth/logout

This API removes the `token` cookie set in the browser, so simply a logout.

#### Call Example:

```
curl -X POST -H 'Content-Type: application/json' -i http://localhost:8080/auth/logout
```

**Output:**

```
"Logged out."
```

---

### GET /dataCollection/status

This API provides the status of the data collector module. It provides the following information:

- **ChannelsRunning**: This is a boolean value indicating that if the channel<sup>[1](#channelFootnote)</sup> extraction process is running at the moment on the server.
- **SensorsRunning**: This is a boolean value indicating that if the sensor/sensor-values extraction process is running at the moment on the server.
- **SensorsProgress**: This is a numeric value indicating the progress of sensor data extraction. It can be from `0` to `100`.
- **NewExtractedChannels**: Once the extraction finishes, this value indicates the number of newly extracted channels.
- **NewExtractedSensors**: Once the extraction finishes, this value indicates the number of newly extracted sensors.
- **NewExtractedSensorValues**:Once the extraction finishes, this value indicates the number of newly extracted sensor values (readings).
- **LastExtractionTime**: This is obvious.

<a name="channelFootnote">1</a>:: We consider a `channel` in ThingSpeak as a `device` in the simulator where can have multiple `sensors` attached to it. So in the API definition and in the database, we keep the ThingSpeak terminology, but in the UI for comfort of the user, we use Waziup terminology.

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -i http://localhost:8080/dataCollection/status
```

**Output:**

```
{
  "ChannelsRunning": false,
  "SensorsRunning": false,
  "SensorsProgress": 100,
  "NewExtractedChannels": 0,
  "NewExtractedSensors": 0,
  "NewExtractedSensorValues": 5144,
  "LastExtractionTime": "2021-06-10T11:22:22.671568363Z"
}
```

---

### GET /dataCollection/statistics

This API provides statistical information on the collected data. It provides the follwoin infomation:

- **totalChannels**: Total number of channels extracted so far.
- **totalSensors**: Total number of sensors extracted so far.
- **totalSensorValues**: Total number of sensor values (readings) extracted so far.

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -i http://localhost:8080/dataCollection/statistics
```

**Output:**

```
{
  "totalChannels": 2450,
  "totalSensorValues": 264596,
  "totalSensors": 4076
}
```

---

### GET /sensors

This API retrieves the information of all sensors.

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -i http://localhost:8080/sensors
```

**Output:**

```
{
  "pagination": {
    "current_page": 1,
    "total_entries": 4076,
    "total_pages": 21
  },
  "rows": [
    {
      "channel_id": 5683,
      "channel_name": "Residential Data Points",
      "id": 393,
      "name": "hot wtr heater"
    },
    {
      "channel_id": 5683,
      "channel_name": "Residential Data Points",
      "id": 401,
      "name": "duct temp"
    },
    {
      "channel_id": 5683,
      "channel_name": "Residential Data Points",
      "id": 409,
      "name": "solar inverter"
    },
    ...
  ]
}
```

---

### GET /sensors/:sensor_id

This API retrieves the information of a sensor for which the `id` is provided.

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -i http://localhost:8080/sensors/409
```

**Output:**

```
{
  "channel_id": 5683,
  "id": 409,
  "name": "solar inverter"
}
```

---

### GET /sensors/:sensor_id/values

This API retrieves all the values of a sensor for which the `id` is provided.

**NOTE**: This API is equivalent to `/channels/:channel_id/sensors/:sensor_id/values`

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -i http://localhost:8080/sensors/409/values
```

**Output:**

```
{
  "pagination": {
    "current_page": 1,
    "total_entries": 58,
    "total_pages": 1
  },
  "rows": [
    {
      "created_at": "2021-06-10T11:16:31Z",
      "entry_id": 1155,
      "name": "solar inverter",
      "sensor_id": 409,
      "value": "13"
    },
    {
      "created_at": "2021-06-10T11:15:33Z",
      "entry_id": 1153,
      "name": "solar inverter",
      "sensor_id": 409,
      "value": "13"
    },
    {
      "created_at": "2021-06-10T11:15:02Z",
      "entry_id": 1152,
      "name": "solar inverter",
      "sensor_id": 409,
      "value": "13"
    },
    ...
  ]
}
```

---

### GET /sensors/:sensor_id/pushSettings [auth required]

This API retrieves all the push settings that are set for a sensor for which the `id` is provided.

_Note: This API requires an authorization token._

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -H 'Authorization: Bearer $2a$10$45Fxw8RvDTT7nLspVKIt9eEna6j0s50dHKjmJDgp0oeRTodPKQeu2' -i http://localhost:8080/sensors/350/pushSettings
```

**Output:**

```
{
  "pagination": {
    "current_page": 1,
    "total_entries": 1,
    "total_pages": 1
  },
  "rows": [
    {
      "active": true,
      "id": 1,
      "last_push_time": null,
      "push_interval": 5,
      "pushed_count": 0,
      "target_device_id": "_49",
      "target_sensor_id": "BAT",
      "use_original_time": true
    }
  ]
}
```

---

### POST /sensors/:sensor_id/pushSettings [auth required]

This API configures a new push setting for a sensor for which the `id` is provided.

_Note: This API requires an authorization token._

#### Input Format:

```
{
  "target_device_id": <String>,
  "target_sensor_id": <String>,
  "active": <Boolean>,
  "push_interval": <Number>,
  "use_original_time": <Boolean>
}
```

#### Call Example:

```
curl -X POST -H 'Content-Type: application/json' -H 'Authorization: Bearer $2a$10$bBPJqUbsTpw9UhirJ.RmIeByMDEstmWAHSQWp.FR19N4aZtptRxBC' -i http://localhost:8080/sensors/350/pushSettings --data '{"target_device_id": "_49","target_sensor_id": "BAT","active": true,"push_interval": 10,"use_original_time": false}'
```

**Output:**

```
OK
```

---

### DELETE /sensors/:sensor_id/pushSettings/:id [auth required]

This API removes the push setting with the given `id` for the determined `sensor_id`.

_Note: This API requires an authorization token._

#### Call Example:

```
curl -X DELETE -H 'Content-Type: application/json' -H 'Authorization: Bearer $2a$10$bBPJqUbsTpw9UhirJ.RmIeByMDEstmWAHSQWp.FR19N4aZtptRxBC' -i http://localhost:8080/sensors/350/pushSettings/2
```

**Output:**

```
OK
```

---

### GET /myPushSettings/sensors [auth required]

This API retrieves all the sensors that the authorized user has set at least a push setting for.

_Note: This API requires an authorization token._

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -H 'Authorization: Bearer $2a$10$45Fxw8RvDTT7nLspVKIt9eEna6j0s50dHKjmJDgp0oeRTodPKQeu2' -i http://localhost:8080/myPushSettings/sensors
```

**Output:**

```
{
  "pagination": {
    "current_page": 1,
    "total_entries": 1,
    "total_pages": 1
  },
  "rows": [
    {
      "channel_id": 215639,
      "id": 350,
      "name": "Solarwatts"
    }
  ]
}
```

---

### GET /search/sensors/:query

This API searches through the collected sensors and retrieves the matching sensors.

_Note: This API requires an authorization token._

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -i http://localhost:8080/search/sensors/temp
```

**Output:**

```
{
  "pagination": {
    "current_page": 1,
    "total_entries": 749,
    "total_pages": 4
  },
  "query": "temp",
  "rows": [
    {
      "channel_id": 1411589,
      "channel_name": "IOT SUMMER TRAINING 2021-PROJECT-1",
      "id": 70,
      "name": "Temp"
    },
    {
      "channel_id": 1411523,
      "channel_name": "IoT OlehV",
      "id": 90,
      "name": "temperature"
    },
    {
      "channel_id": 1411292,
      "channel_name": "DHT11",
      "id": 99,
      "name": "Temp (C°)"
    },
    {
      "channel_id": 1409789,
      "channel_name": "temperature test",
      "id": 814,
      "name": "temperature"
    },
    ...
  ]
}
```

---

### GET /channels

This API retrieves all the channels.

**Note**: We consider a `channel` in ThingSpeak as a `device` in the simulator where can have multiple `sensors` attached to it. So in the API definition and in the database, we keep the ThingSpeak terminology, but in the UI for comfort of the user, we use Waziup terminology.

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -i http://localhost:8080/channels
```

**Output:**

```
{
  "pagination": {
    "current_page": 1,
    "total_entries": 2488,
    "total_pages": 13
  },
  "rows": [
    {
      "created_at": "2021-04-27T07:02:18Z",
      "description": "VEHICLE DATA",
      "id": 1372826,
      "last_entry_id": 0,
      "latitude": 0,
      "longitude": 0,
      "name": "IDA",
      "url": ""
    },
    {
      "created_at": "2021-05-08T13:55:41Z",
      "description": "Iot based air quality monitoring system with alert",
      "id": 1384419,
      "last_entry_id": 0,
      "latitude": 0,
      "longitude": 0,
      "name": "Project_s8_air_quality",
      "url": ""
    },
    ...
  ]
}
```

---

### GET /channels/:channel_id

This API retrieves details of a channel with the given `id`.

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -i http://localhost:8080/channels/1402239
```

**Output:**

```
{
  "created_at": "2021-05-28T09:48:13Z",
  "description": "board with real sensors data",
  "id": 1402239,
  "last_entry_id": 0,
  "latitude": 45.791964,
  "longitude": 15.961711,
  "name": "real_board",
  "url": ""
}
```

---

### GET /channels/:channel_id/sensors

This API retrieves the sensors and the details of a channel with the given `channel_id`.

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -i http://localhost:8080/channels/215639/sensors
```

**Output:**

```
{
  "channel": {
    "created_at": "2017-01-18T11:07:58Z",
    "description": "Temperature from Current cost via HAH.",
    "id": 215639,
    "last_entry_id": 3558484,
    "latitude": 0,
    "longitude": 0,
    "name": "House Monitor",
    "url": ""
  },
  "pagination": {
    "current_page": 1,
    "total_entries": 8,
    "total_pages": 1
  },
  "rows": [
    {
      "id": 349,
      "name": "Temp"
    },
    {
      "id": 350,
      "name": "Solarwatts"
    },
    {
      "id": 351,
      "name": "Gas"
    },
    {
      "id": 352,
      "name": "ConsumedW"
    },
    {
      "id": 353,
      "name": "SpareC"
    },
    {
      "id": 354,
      "name": "Flow"
    },
    {
      "id": 355,
      "name": "Return"
    },
    {
      "id": 356,
      "name": "Solar Export"
    }
  ]
}
```

---

### GET /user

This API retrieves the details of the authorized user.

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -H 'Authorization: Bearer $2a$10$M35K.pfQs6714ZArnE1MMOlw.Rd.E84.enCeRyeSxtaS22XHPk03q' -i http://localhost:8080/user
```

**Output:**

```
{
  "tokenHash": "$2a$10$M35K.pfQs6714ZArnE1MMOlw.Rd.E84.enCeRyeSxtaS22XHPk03q",
  "username": "cdupont"
}
```

---

### GET /userDevices

This API retrieves the devices of the user from Waziup cloud. It returns the Waziup API directly.

#### Call Example:

```
curl -X GET -H 'Content-Type: application/json' -H 'Authorization: Bearer $2a$10$M35K.pfQs6714ZArnE1MMOlw.Rd.E84.enCeRyeSxtaS22XHPk03q' -i http://localhost:8080/userDevices
```

**Output:**

```
[
  {
    "gateway_id": "",
    "date_modified": "2021-05-26T17:49:05Z",
    "domain": "Ethiopia",
    "visibility": "public",
    "owner": "cdupont",
    "name": "SensoDyn",
    "id": "SM2222",
    "sensors": [...],
    "actuators": [...],
    "date_created": "2019-06-03T16:14:01Z"
  },
  {
    "date_modified": "2021-04-01T08:19:56Z",
    "location": {
      "latitude": 43.517200698,
      "longitude": 8.833007813
    },
    "domain": "waziup",
    "visibility": "public",
    "owner": "cdupont",
    "name": "My device 89 b",
    "id": "MyDevice89",
    "sensors": [...],
    "actuators": [...],
    "date_created": "2019-06-07T06:52:25Z"
  },
  ...
]
```

---
