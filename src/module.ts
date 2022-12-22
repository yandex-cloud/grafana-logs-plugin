import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { LoggingQuery, LoggingSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, LoggingQuery, LoggingSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
