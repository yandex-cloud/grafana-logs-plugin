import { DataQuery, DataSourceJsonData } from '@grafana/data';

export const knownLevels: string[] = [
  "TRACE",
  "DEBUG",
  "INFO",
  "WARN",
  "ERROR",
  "FATAL",
];


export interface LoggingQuery extends DataQuery {
  groupId: string;

  limit: number;
  queryText?: string;
  levels?: string[];
  stream?: string;
  resourceType?: string;
  resourceIds?: string[];
  addPayloadFields?: string[];
}

export const defaultQuery: Partial<LoggingQuery> = {
  limit: 10,
  levels: knownLevels,
};

/**
 * These are options configured for each DataSource instance
 */
export interface LoggingSourceOptions extends DataSourceJsonData {
  apiEndpoint: string;
  folderId: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface LoggingSecureJsonData {
  apiKeyJson?: string;
}

export const defaultSourceOptions: Partial<LoggingSourceOptions> = {
  apiEndpoint: "api.cloud.yandex.net:443",
  folderId: "",
}
