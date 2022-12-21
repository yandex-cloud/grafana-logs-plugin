import defaults from 'lodash/defaults';

import React, { ChangeEvent, PureComponent } from 'react';
import {
  Field, TextArea, Input, InlineFieldRow, InlineField, MultiSelect, Select, TagsInput
} from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, knownLevels, LoggingSourceOptions, LoggingQuery } from './types';


type Props = QueryEditorProps<DataSource, LoggingQuery, LoggingSourceOptions>;

interface QueryEditorCache {
  groups?: string[];
  resourceTypes?: string[];
  resourceIds?: string[];
}

interface QueryCacheParams {
  forGroup?: string
  forResourceType?: string
}

interface QueryEditorSuggestRequest {
  groupId: string;
  resourceType: string;
}

type State = {
  cache: QueryEditorCache;
  cachedFor?: QueryCacheParams;
}

export class QueryEditor extends PureComponent<Props, State> {

  state: Readonly<State> = { cache: {} }

  updateCache = async () => {
    const { query, datasource } = this.props
    const { cache, cachedFor } = this.state
    if ((cachedFor !== undefined) && query.groupId === cachedFor.forGroup && query.resourceType === cachedFor.forResourceType) {
      return
    }

    const req: QueryEditorSuggestRequest = {
      groupId: query.groupId || "",
      resourceType: query.resourceType || "",
    }

    await datasource.postResource("suggestQuery", req).then((v) => {
      this.setState({
        cachedFor: {
          forGroup: query.groupId,
          forResourceType: query.resourceType,
        },
        cache: {
          ...cache, ...v,
        },
      });
    })
  };

  onGroupIdChange = (event: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, groupId: event.value || "" });
    this.runQueryIfNeeded();
  };

  onLimitChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, limit: parseFloat(event.target.value) });
    this.runQueryIfNeeded();
  };

  onLevelsChange = (events: Array<SelectableValue<string>>) => {
    const { onChange, query } = this.props;
    let chosen = events.map((op): string => (op.value!))
    if (chosen.length < 1) {
      chosen = knownLevels
    }
    onChange({ ...query, levels: chosen });
    this.runQueryIfNeeded();
  };

  onQueryTextChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, queryText: event.target.value });
  };

  onStreamChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, stream: event.target.value });
    this.runQueryIfNeeded();
  };

  onResourceTypeChange = (event: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, resourceType: event.value });
    this.runQueryIfNeeded();
  };

  onResourceIdsChange = (events: Array<SelectableValue<string>>) => {
    const { onChange, query } = this.props;
    onChange({ ...query, resourceIds: events.map((op): string => (op.value!)) });
    this.runQueryIfNeeded();
  };

  onAddPayloadFieldsChange = (values: string[]) => {
    const { onChange, query } = this.props;
    onChange({ ...query, addPayloadFields: values });
    this.runQueryIfNeeded();
  };

  runQueryIfNeeded = () => {
    const { query, onRunQuery } = this.props;
    console.log(query)
    if (!query.groupId) {
      return;
    }
    if ((query.levels?.length ?? 0) < 1) {
      return
    }
    onRunQuery();
  };

  render() {
    this.updateCache()


    const query = defaults(this.props.query, defaultQuery);
    const { cache } = this.state
    const { groupId, limit, levels, stream, resourceIds, resourceType, queryText, addPayloadFields } = query;

    const levelsOptions: Array<SelectableValue<string>> = knownLevels.map((l): SelectableValue<string> => ({ value: l, label: l }));
    const groupOptions: Array<SelectableValue<string>> = [
      ...(groupId !== "" ? [groupId] : []),
      ...(cache.groups || [])
    ].filter((v, i, a) => a.indexOf(v) === i).map((l): SelectableValue<string> => ({ value: l, label: l }))
    const resourceTypeOptions: Array<SelectableValue<string>> = [
      "", // no resource type
      ...(resourceType !== "" ? [resourceType] : []),
      ...(cache.resourceTypes || [])
    ].filter((v, i, a) => a.indexOf(v) === i).map((l): SelectableValue<string> => ({ value: l, label: l }))
    const resourceIdsOptions: Array<SelectableValue<string>> = [
      ...(resourceIds || []), ...(cache.resourceIds || [])
    ].filter((v, i, a) => a.indexOf(v) === i).map((l): SelectableValue<string> => ({ value: l, label: l }))

    return (
      <div>
        <InlineFieldRow>
          <InlineField label='Group'>
            <Select
              allowCustomValue={true}
              value={groupId}
              options={groupOptions}
              onChange={this.onGroupIdChange}
            />
          </InlineField>
          <InlineField label='Limit'>
            <Input
              type='number'
              name='log-limit'
              value={limit}
              min={1}
              onChange={this.onLimitChange}
            />
          </InlineField>
          <InlineField label='Levels'>
            <MultiSelect
              options={levelsOptions}
              value={levels}
              closeMenuOnSelect={false}
              onChange={this.onLevelsChange}
            />
          </InlineField>
        </InlineFieldRow>
        <InlineFieldRow>
          <InlineField label='Stream'>
            <Input
              type='text'
              name='log-stream'
              value={stream}
              onChange={this.onStreamChange}
            />
          </InlineField>
          <InlineField label='Resource type'>
            <Select
              allowCustomValue={true}
              value={resourceType}
              options={resourceTypeOptions}
              onChange={this.onResourceTypeChange}
            />
          </InlineField>
          <InlineField label='Resource id'>
            <MultiSelect
              allowCustomValue={true}
              value={resourceIds}
              options={resourceIdsOptions}
              onChange={this.onResourceIdsChange}
            />
          </InlineField>
        </InlineFieldRow>
        <Field label='Filter query' description='Query to filter log records'>
          <TextArea
            name='query-text'
            value={queryText}
            onChange={this.onQueryTextChange}
          />
        </Field>
        <InlineField label="Add payload fields to content" >
          <TagsInput
            onChange={this.onAddPayloadFieldsChange}
            tags={addPayloadFields}
          />
        </InlineField>
      </div >
    );
  };
}
