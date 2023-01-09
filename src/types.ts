import { DataQuery, DataSourceJsonData } from '@grafana/data';

export const knownLevels: string[] = [
  "TRACE",
  "DEBUG",
  "INFO",
  "WARN",
  "ERROR",
  "FATAL",
];

export interface DerivedFiledRule {
  name: string;
  template: string;
}

export interface LoggingQuery extends DataQuery {
  groupId: string;

  limit: number;
  queryText?: string;
  levels?: string[];
  stream?: string;
  resourceType?: string;
  resourceIds?: string[];
  addPayloadFields?: string[];
  derivedFields?: DerivedFiledRule[];
}

export const defaultQuery: Partial<LoggingQuery> = {
  limit: 10,
  levels: knownLevels,
};


export interface LoggingDerivedLink {
  field: string;
  title: string;
  url: string;
  targetBlank: boolean;
}

/**
 * These are options configured for each DataSource instance
 */
export interface LoggingSourceOptions extends DataSourceJsonData {
  apiEndpoint: string;
  folderId: string;
  derivedLinks?: LoggingDerivedLink[];
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
