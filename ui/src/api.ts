const API_URL = "/";

/*--------------*/

export type HttpError = {
  status: number;
  message: string;
};

async function failResp(resp: Response) {
  var text = await resp.text();
  throw { status: resp.status, message: text } as HttpError;
}

/*--------------*/

export type DataCollectionStatus = {
  ChannelsRunning: boolean;
  SensorsRunning: boolean;
  SensorsProgress: number;
  NewExtractedChannels: number;
  NewExtractedSensors: number;
  NewExtractedSensorValues: number;
  LastExtractionTime: Date;
};

export async function getDataCollectionStatus(): Promise<DataCollectionStatus> {
  var resp = await fetch(API_URL + "dataCollection/status");
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export type DataCollectionStatistics = {
  totalChannels: number;
  totalSensors: number;
  totalSensorValues: number;
};

export async function getDataCollectionStatistics(): Promise<DataCollectionStatistics> {
  var resp = await fetch(API_URL + "dataCollection/statistics");
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export type Pagination = {
  current_page: number;
  total_entries: number;
  total_pages: number;
};

export type ChannelRow = {
  created_at: Date;
  description: string;
  id: number;
  last_entry_id: number;
  latitude: number;
  longitude: number;
  name: string;
  url: string;
};

export type ChannelsData = {
  pagination: Pagination;
  rows: ChannelRow[];
};

export async function getChannels(page: number): Promise<ChannelsData> {
  if (!page) {
    page = 1;
  }
  var resp = await fetch(API_URL + "channels?page=" + page.toString());
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export async function getChannel(id: number): Promise<ChannelRow> {
  var resp = await fetch(`${API_URL}channels/${id.toString()}`);
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export type SensorRow = {
  id: number;
  name: string;
  channel_id: number;
  channel_name: string;
};

export type SensorsData = {
  pagination: Pagination;
  rows: SensorRow[];
};

export async function searchSensors(
  query: string,
  page: number
): Promise<SensorsData> {
  if (!page) {
    page = 1;
  }
  query = encodeURI(query);
  var resp = await fetch(
    `${API_URL}search/sensors/${query}?page=${page.toString()}`
  );
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export type SensorValueRow = {
  entry_id: number;
  sensor_id: number;
  name: string;
  value: string;
  created_at: Date;
};

export type SensorValues = {
  pagination: Pagination;
  rows: SensorValueRow[];
};

export async function getSensorValues(
  id: number,
  page: number
): Promise<SensorValues> {
  if (!page) {
    page = 1;
  }
  var resp = await fetch(
    `${API_URL}sensors/${id.toString()}/values?page=${page.toString()}`
  );
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export async function getSensor(id: number): Promise<SensorRow> {
  var resp = await fetch(`${API_URL}sensors/${id.toString()}`);
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export type UserType = {
  username: string;
  tokenHash: string;
};

export async function login(
  username: string,
  password: string
): Promise<string> {
  var resp = await fetch(`${API_URL}auth`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });
  if (!resp.ok) await failResp(resp);
  return await resp.text();
  // return await resp.json();
}

/*--------------*/

export async function logout(): Promise<string> {
  var resp = await fetch(`${API_URL}auth/logout`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    // body: JSON.stringify({}),
  });
  if (!resp.ok) await failResp(resp);
  return await resp.text();
  // return await resp.json();
}

/*--------------*/

export async function getUser(): Promise<UserType> {
  var resp = await fetch(`${API_URL}user`); // The tokenHash will be sent through cookies
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export async function getUserDevices(): Promise<any> {
  var resp = await fetch(`${API_URL}userDevices`); // The tokenHash will be sent through cookies
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export type SensorPushSettings = {
  id?: number;
  target_device_id: string;
  target_sensor_id: string;
  active: boolean;
  push_interval: number;
  last_push_time?: Date;
};

export type AllSensorPushSettings = {
  pagination: Pagination;
  rows: SensorPushSettings[];
};

export async function getPushSettings(
  sensor_id: number,
  page: number
): Promise<AllSensorPushSettings> {
  if (!page) {
    page = 1;
  }
  var resp = await fetch(
    `${API_URL}sensors/${sensor_id}/pushSettings?page=${page.toString()}`
  ); // The tokenHash will be sent through cookies
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

export async function savePushSettings(
  record: SensorPushSettings,
  sensor_id: number
): Promise<any> {
  var resp = await fetch(`${API_URL}sensors/${sensor_id}/pushSettings`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(record),
  });
  if (!resp.ok) await failResp(resp);
  return await resp.text();
  // return await resp.json();
}

/*--------------*/

export async function deletePushSettings(
  sensor_id: number,
  recordId: number
): Promise<any> {
  var resp = await fetch(
    `${API_URL}sensors/${sensor_id}/pushSettings/${recordId}`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    }
  );
  if (!resp.ok) await failResp(resp);
  return await resp.text();
  // return await resp.json();
}

/*--------------*/

export type SensorType = {
  id: string;
  name: string;
  devName: string;
  devId: string;
  title: string; // devName + name
};

/*--------------*/
