import defaults from 'lodash/defaults';

import {
  DataSourceInstanceSettings,
  ScopedVars,
} from '@grafana/data';

import {
  DataSourceWithBackend,
  getTemplateSrv
} from '@grafana/runtime';

import { LoggingQuery, LoggingSourceOptions, defaultQuery } from './types';

export class DataSource extends DataSourceWithBackend<LoggingQuery, LoggingSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<LoggingSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(inQuery: LoggingQuery, scopedVars: ScopedVars): Record<string, any> {
    const tsrv = getTemplateSrv();
    const query = defaults(inQuery, defaultQuery);

    query.groupId = tsrv.replace(query.groupId, scopedVars);
    query.queryText = tsrv.replace(query.queryText, scopedVars);
    query.resourceType = tsrv.replace(query.resourceType, scopedVars);
    if ((query.resourceIds || []).length === 1) {
      const resReplaced = tsrv.replace(query.resourceIds![0], scopedVars, "json");
      try {
        const resList = JSON.parse(resReplaced);
        query.resourceIds = resList;
      } catch {
        // skip var replacement
      };
    };
    query.stream = tsrv.replace(query.stream, scopedVars);
    return query;
  };
}
