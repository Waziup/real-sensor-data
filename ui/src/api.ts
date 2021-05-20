const LOCAL_API_URL = "/";
const WAZIUP_API_URL = "....";

/*--------------*/

async function failResp(resp: Response) {
  var text = await resp.text();
  throw `There was an error calling the API.\nThe server returned (${resp.status}) ${resp.statusText}.\n\n${text}`;
}

/*--------------*/

export type DataCollectionStatus = {
  ChannelsRunning: boolean;
  SensorsRunning: boolean;
  SensorsProgress: number;
  NewExtractedChannels: number;
  NewExtractedSensorValues: number;
  LastExtractionTime: Date;
};

export async function getDataCollectionStatus(): Promise<DataCollectionStatus> {
  var resp = await fetch(LOCAL_API_URL + "dataCollection/status");
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
  var resp = await fetch(LOCAL_API_URL + "dataCollection/statistics");
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
  var resp = await fetch(LOCAL_API_URL + "channels?page=" + page.toString());
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export async function getChannel(id: number): Promise<ChannelRow> {
  var resp = await fetch(`${LOCAL_API_URL}channels/${id.toString()}`);
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export type SensorRow = {
  channel_id: number;
  channel_name: string;
  name: string;
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
    `${LOCAL_API_URL}sensors/search/${query}?page=${page.toString()}`
  );
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/

export type SensorValueRow = {
  channel_id: number;
  created_at: Date;
  entry_id: number;
  name: string;
  value: string;
};

export type SensorValues = {
  pagination: Pagination;
  rows: SensorValueRow[];
};

export async function getSensorValues(
  name: string,
  channel_id: number
): Promise<SensorValues> {
  name = encodeURI(name);
  var resp = await fetch(
    `${LOCAL_API_URL}channels/${channel_id.toString()}/sensors/${name}/values`
  );
  if (!resp.ok) await failResp(resp);
  return await resp.json();
}

/*--------------*/
