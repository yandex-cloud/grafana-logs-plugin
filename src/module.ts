import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { LoggingQuery, LoggingSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, LoggingQuery, LoggingSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
